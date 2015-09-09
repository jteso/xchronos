package agent

import (
	"fmt"
	"log"
	"time"

	"github.com/jteso/task"
	"github.com/jteso/xchronos/cluster"
	"github.com/jteso/xchronos/scheduler"
)

type Agent struct {
	// agent's id
	ID string
	// internal ip
	IPv4 string
	// Agent's state
	state string
	// Last error reported by the agent, or agent's task
	lastError error
	// Used to communicate to etcd cluster
	clusterClient cluster.ClusterClient
	// Manager of the all tasks running on the background by an agent
	taskManager *task.TaskManager
	// Debugging flag
	verbose bool
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
		state:         "INIT",
		clusterClient: cluster.NewEtcdClient(etcdNodes),
		verbose:       verbose,
		haltTaskC:     make(chan struct{}),
		taskManager:   task.NewTaskManager(),
	}
}

func (a *Agent) Run() error {
	for stateHandler := startStateFn; stateHandler != nil; {
		stateHandler = stateHandler(a)
	}
	//close(a.stopCh)
	return nil
}

func (a *Agent) runForLeader() (bool, error) {
	return a.clusterClient.SchedulerElect(a.IPv4)
}

func (a *Agent) connectEtcdCluster() error {
	return a.clusterClient.Connect()
}

func (a *Agent) disconnectEtcdCluster() error {
	return a.clusterClient.Disconnect()
}

func (a *Agent) changeState(newState string) {
	a.log(fmt.Sprintf("Changing state: %s -> %s", a.state, newState))
	a.state = newState
}

func (a *Agent) advertiseAndRenewLeaderRoleT() *task.Task {
	key := SCHEDULER_ELECTION_KEY
	t := task.New("leaderRenewal", func() error {
		a.log("Renewing my leader role...")
		_, err := a.clusterClient.SchedulerElect(a.IPv4)
		return err
	})
	t.RunEvery(time.Second * HEARTBEAT)
	a.taskManager.RegisterTask(t)
	return t
}

// takeExecutorRole function will make the agent to offer itself to execute jobs been offered.
// this offering is been done by writing periodically (HEARTBEAT) into the etcd dir /executors
// this function will returned via chan any error is encountered, and this agent can be stop been
// offered as an executor by stopping the agents executorTicker.
func (a *Agent) advertiseAndRenewExecutorRoleT() *task.Task {
	t := task.New("executorRenewal", func() error {
		a.log("Renewing my executor role...")
		_, err := a.clusterClient.RegisterAsExecutor(a.ID, a.IPv4)
		return err
	})
	t.RunEvery(time.Second * HEARTBEAT)
	a.taskManager.RegisterTask(t)
	return t
}

func (a *Agent) publishJobOffersT(job *Job) *task.Task {
	t := task.New("jobOffersPublisher", func() error {
		return a.clusterClient.MakeJobOffer(job)
	})

	t.RunEvery(time.Second * 1)
	a.taskManager.RegisterTask(t)
	return t
}

func (a *Agent) watchForJobOffersT() *task.Task {
	jobOfferC := make(chan *scheduler.Job, 1)
	jobOfferStopC := make(chan bool, 1)
	t := task.New("jobOffersWatcher", func() error {
		a.clusterClient.WatchJobOffers(jobOfferC, jobOfferStopC)
		var jobOffer *scheduler.Job
		for {
			jobOffer = <-jobOfferC
			if a.clusterClient.TakeJobOffer(jobOffer, a.IPv4) {
				a.logf("Job offer accepted by agent: [%s]", a.IPv4)
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

func (a *Agent) watchForNewLeaderElectionT() *task.Task {
	notifyC := make(chan bool, 1)
	watchLeaderStopC := make(chan bool, 1)

	t := task.New("watchLeaderElection", func() error {
		a.clusterClient.NotifySchedulerStepDown(notifyC, watchLeaderStopC)
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

func (a *Agent) Stop() {
	doneC := make(chan bool, 1)
	a.taskManager.StopAllTasks(doneC)
	<-doneC
	a.logf("Agent halted")
}

// func (a *Agent) registerTask(newTask *task.Task) {
// 	if len(a.taskManager) == 0 {
// 		// internal to the task manager that captures cancelation signals
// 		// from users and trigger the stop of all tasks
// 		t := task.New("uiTask", func() error {
// 			<-a.haltTaskC
// 			return task.ErrUserCanceled
// 		})
// 		t.RunOnce()
// 		a.taskManager = append(a.taskManager, t)
// 	}
// 	if a.verbose {
// 		a.logf("New Task registered: %s", newTask.Id)
// 	}
// 	a.taskManager = append(a.taskManager, newTask)
// }

// func (a *Agent) listenUICancelTask() *task.Task {
// 	return a.taskManager[0]
// }

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
