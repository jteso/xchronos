package agent

import (
	"fmt"
	"log"
	"time"

	"github.com/jteso/task"
	"github.com/jteso/xchronos/cluster"
	"github.com/jteso/xchronos/errors"
	"github.com/jteso/xchronos/scheduler"
	"github.com/jteso/xchronos/supervisor"
)

type AgentSpec interface {
	// Agent gets notified when it has transition into a new state
	NotifyState(newState string)
	// Connect into the job store
	ConnectCluster()
	// Disconnect from the jobstore (etcd, consul,...)
	DisconnectCluster()
	// This function is used for leader election, been the leader the
	// agent responsible to manage the job scheduler
	AttemptJobScheduler() (bool, error)
	// Agent will be responsible to manage the job scheduler
	ActAsJobScheduler() chan error
	// Agent will be only responsible to accept job offers, and attempt their execution
	ActAsJobExecutor() chan error
	// Stop all tasks that an agent may be executing as part of this role (scheduler or executor)
	Halt()
}

type Agent struct {
	// agent's id
	ID string
	// internal ip
	IPv4 string
	// Agent's state
	state string
	// Used to communicate to etcd cluster
	clusterClient cluster.ClusterClient
	// Manager of the all tasks running on the background by an agent
	taskManager *task.TaskManager
	// Last error reported by the agent, or agent's task
	lastError error
	// internal scheduler
	jobScheduler *scheduler.Scheduler
	// Debugging flag
	verbose bool
	// Build-in supervisor to run the jobs
	jobSupervisor *supervisor.Supervisor
}

// etcdNotes = strings.Split(os.Getenv("ETCD_NOTES"), ",")
func New(id string, etcdNodes []string, verbose bool) *Agent {
	localIp, err := GetLocalIPv4()
	if err != nil {
		panic(err)
	}

	return &Agent{
		ID:            id,
		IPv4:          localIp,
		state:         "_init",
		clusterClient: cluster.NewEtcdClient(etcdNodes),
		taskManager:   task.NewTaskManager(),
		jobScheduler:  scheduler.NewScheduler(),
		verbose:       verbose,
		jobSupervisor: new(supervisor.Supervisor),
	}
}

func (a *Agent) Run() error {
	for stateHandler := startStateFn; stateHandler != nil; {
		stateHandler = stateHandler(a)
	}
	return nil
}

func (a *Agent) NotifyState(newState string) {
	a.logf("Changing state: %s -> %s", a.state, newState)
	a.state = newState
}

func (a *Agent) ConnectCluster() error {
	return a.clusterClient.Connect()
}

func (a *Agent) DisconnectCluster() {
	a.clusterClient.Disconnect()
}

func (a *Agent) AttemptJobScheduler() (bool, error) {
	return a.clusterClient.RegisterAsScheduler(a.IPv4)
}

func (a *Agent) ActAsJobScheduler() chan error {
	errorC := make(chan error, 1)
	// tasks as scheduler
	schedulerRoleTask := a.advertiseSchedulerRoleTask()
	watchForJobsToScheduleTask := a.watchForJobsToScheduleTask()
	publishJobOffersTask := a.publishJobOffersTask()

	// tasks as job executor
	executorRoleTask := a.advertiseExecutorRoleTask()
	offersListenerTask := a.watchForJobOffersTask()
	schedulerFailureWatcherTask := a.watchForSchedulerFailureTask()

	go func() {
		for {
			select {
			case err := <-task.FirstError(
				schedulerRoleTask,
				watchForJobsToScheduleTask,
				publishJobOffersTask,
				executorRoleTask,
				offersListenerTask):

				a.lastError = err
				errorC <- err
			case <-schedulerFailureWatcherTask.ErrorChan():
				a.lastError = errors.ErrNoSchedulerDetected
				errorC <- errors.ErrNoSchedulerDetected
			}
		}
	}()
	return errorC
}

func (a *Agent) ActAsJobExecutor() chan error {
	errorC := make(chan error, 1)

	// tasks as job executor
	executorRoleTask := a.advertiseExecutorRoleTask()
	offersListenerTask := a.watchForJobOffersTask()
	schedulerFailureWatcherTask := a.watchForSchedulerFailureTask()

	go func() {
		for {
			select {
			case err := <-task.FirstError(
				executorRoleTask,
				offersListenerTask):

				a.lastError = err
				errorC <- err
			case <-schedulerFailureWatcherTask.ErrorChan():
				a.lastError = errors.ErrNoSchedulerDetected
				errorC <- errors.ErrNoSchedulerDetected
			}
		}
	}()
	return errorC
}

// This function register the leader (aka scheduler) IP periodically.
// Keep in mind that the "title" of leader expires every `SCHEDULER_LEADER_TTL`
// and it has to be renewed. Fail to do that, the rest of cluster will compete
// again for a leaders position.
func (a *Agent) advertiseSchedulerRoleTask() *task.Task {
	t := task.New("advertiseSchedulerRoleTask", func() error {
		a.log("Renewing my leader role...")
		_, err := a.clusterClient.RegisterAsScheduler(a.IPv4)
		return err
	})
	t.RunEvery(time.Second * cluster.HEARTBEAT)
	a.taskManager.RegisterTask(t)
	return t
}

// This function listen for new jobs pending for scheduling
func (a *Agent) watchForJobsToScheduleTask() *task.Task {
	unscheduledJobsC := make(chan *scheduler.Job, 100)
	stopDueJobC := make(chan bool, 1)

	t := task.New("watchForJobsToScheduleTask", func() error {
		go a.clusterClient.WatchJobsToSchedule(unscheduledJobsC, stopDueJobC)

		a.log("Job sentinel activated")
		for unscheduledJob := range unscheduledJobsC {
			if unscheduledJob != nil {
				a.logf("Scheduling job: %s (due in %f)...", unscheduledJob.ToString(), unscheduledJob.WaitSecs())
				a.jobScheduler.Enqueue(unscheduledJob)
			}
		}
		return nil
	})
	t.OnStopFn(func() {
		stopDueJobC <- true
	})
	t.RunOnce()
	a.taskManager.RegisterTask(t)
	return t
}

// This function will be run by the leader only, and it will manage the job scheduler
// and by extent, the publication of job offers
func (a *Agent) publishJobOffersTask() *task.Task {
	dueJobC := make(chan *scheduler.Job, 10)
	stopDueJobC := make(chan bool)

	t := task.New("publishJobOffersTask", func() error {
		a.jobScheduler.Notify(dueJobC, stopDueJobC)
		for job := range dueJobC {
			a.logf("Publishing jobOffer: %s ...", job.ToString())
			a.clusterClient.MakeJobOffer(job)
		}
		return nil
	})
	t.OnStopFn(func() {
		stopDueJobC <- true
	})
	t.RunOnce()
	a.taskManager.RegisterTask(t)
	return t
}

// this function will make the agent a candidate to receive job offers.
func (a *Agent) advertiseExecutorRoleTask() *task.Task {
	t := task.New("advertiseExecutorRoleTask", func() error {
		a.log("Renewing my executor role...")
		return a.clusterClient.RegisterAsExecutor(a.ID, a.IPv4)
	})

	t.RunEvery(time.Second * cluster.HEARTBEAT)
	a.taskManager.RegisterTask(t)
	return t
}

func (a *Agent) watchForJobOffersTask() *task.Task {
	jobOfferC := make(chan *scheduler.Job, 1)
	jobOfferStopC := make(chan bool, 1)
	t := task.New("watchForJobOffersTask", func() error {
		go a.clusterClient.WatchJobOffers(jobOfferC, jobOfferStopC)
		for jobOffer := range jobOfferC{
			if jobOffer != nil && a.clusterClient.TakeJobOffer(jobOffer, a.IPv4) {
				a.logf("Executing jobOffer: [%s]", jobOffer.ToString())
				go a.jobSupervisor.RunIt(jobOffer)
			}
		}
		return nil
	})
	t.OnStopFn(func() {
		jobOfferStopC <- true
	})
	t.RunOnce()
	a.taskManager.RegisterTask(t)
	return t
}

func (a *Agent) watchForSchedulerFailureTask() *task.Task {
	notifyC := make(chan bool, 1)
	watchLeaderStopC := make(chan bool, 1)

	t := task.New("watchForSchedulerFailureTask", func() error {
		a.clusterClient.SchedulerFailureWatcher(notifyC, watchLeaderStopC)
		<-notifyC
		return nil
	})

	t.OnStopFn(func() {
		watchLeaderStopC <- true
	})
	t.RunOnce()
	a.taskManager.RegisterTask(t)
	return t
}

func (a *Agent) Halt() {
	a.logf("Halting due to error: %s", a.lastError.Error())
	doneC := make(chan bool, 1)
	a.taskManager.StopAllTasks(doneC)
	<-doneC
	a.logf("Agent halted")
}

func (a *Agent) log(message string) {
	if a.verbose {
		log.Printf("[%s] %s\n", a.ID, message)
	}
}

func (a *Agent) logf(format string, v ...interface{}) {
	if a.verbose {
		a.log(fmt.Sprintf(format, v...))
	}
}
