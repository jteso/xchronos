package agent

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/jteso/xchronos/Godeps/_workspace/src/github.com/coreos/go-etcd/etcd"
	"github.com/jteso/xchronos/task"
)

const (
	SCHEDULER_ELECTION_KEY = "/xchronos/var/scheduler/election"
	EXECUTORS_DIR          = "/xchronos/etc/executors"
	JOBS_DIR               = "/xchronos/etc/jobs"
)

var (
	SCHEDULER_LEADER_TTL uint64        = 10 // max time running without leader
	EXECUTOR_TTL         uint64        = 10 // max time running without a particular executor
	HEARTBEAT            time.Duration = 5
)

type Agent struct {
	// agent's id
	ID string
	// Agent's state
	state string
	// Last error reported by the agent, or agent's task
	lastError error

	// Used to communicate to etcd cluster
	etcdClient *etcd.Client
	// See all the etcd nodes available in the cluster
	etcdNodes []string

	// Manager of the all tasks running on the background by an agent
	taskManager []*task.Task
	// Signal by ui to stop all registered tasks
	haltTaskC chan struct{}
	// Signal back to ui, indicating when all registered tasks have stopped
	haltedC chan struct{}

	// receiving job offers channel pending to be executed
	jobC chan *etcd.Response
	// stop watching for jobs
	jobStopC chan bool

	// Debugging flag
	verbose bool
}

// etcdNotes = strings.Split(os.Getenv("ETCD_NOTES"), ",")
func New(id string, etcdNodes []string, verbose bool) *Agent {
	return &Agent{
		ID:          id,
		state:       "INIT",
		etcdNodes:   etcdNodes,
		verbose:     verbose,
		haltTaskC:   make(chan struct{}),
		taskManager: []*task.Task{},
	}
}

func (a *Agent) Run() error {
	for stateHandler := startStateFn; stateHandler != nil; {
		stateHandler = stateHandler(a)
	}
	//close(a.stopCh)
	return nil
}

// runForLeader function will make an agent to run for scheduler leadership.
// output:
// - true, nil : leader
// - false, nil: supporter
// - false, err: error
func (a *Agent) runForLeader() (bool, error) {
	// Put value if prevExist=false
	_, err := a.etcdClient.Create(SCHEDULER_ELECTION_KEY, a.ID, SCHEDULER_LEADER_TTL)

	if err != nil {
		if etcdError, ok := err.(*etcd.EtcdError); ok {
			// key already exists
			if etcdError.ErrorCode == 105 {
				return false, nil
			}
		}

		return false, err
	}
	return true, nil
}

func (a *Agent) connectEtcdCluster() error {
	a.etcdClient = etcd.NewClient(a.etcdNodes)
	return nil
}

func (a *Agent) disconnectEtcdCluster() error {
	a.etcdClient.Close()
	return nil
}

func (a *Agent) changeState(newState string) {
	a.log(fmt.Sprintf("Changing state: %s -> %s", a.state, newState))
	a.state = newState
}

func (a *Agent) advertiseAndRenewLeaderRoleT() *task.Task {
	key := SCHEDULER_ELECTION_KEY
	t := task.New("leaderRenewal", func() error {
		a.log("Renewing my leader role...")
		_, err := a.etcdClient.Set(key, a.ID, SCHEDULER_LEADER_TTL)
		return err
	})
	t.RunEvery(time.Second * HEARTBEAT)
	a.registerTask(t)
	return t
}

// takeExecutorRole function will make the agent to offer itself to execute jobs been offered.
// this offering is been done by writing periodically (HEARTBEAT) into the etcd dir /executors
// this function will returned via chan any error is encountered, and this agent can be stop been
// offered as an executor by stopping the agents executorTicker.
func (a *Agent) advertiseAndRenewExecutorRoleT() *task.Task {
	t := task.New("executorRenewal", func() error {
		a.log("Renewing my executor role...")
		_, err := a.etcdClient.Set(EXECUTORS_DIR+"/"+a.ID, "up", EXECUTOR_TTL)
		return err
	})
	t.RunEvery(time.Second * HEARTBEAT)
	a.registerTask(t)
	return t
}

func (a *Agent) publishJobOffersT() *task.Task {
	jobId := 0
	t := task.New("jobOffersPublisher", func() error {
		key := JOBS_DIR + "/agent_1/" + strconv.Itoa(jobId)
		_, err := a.etcdClient.Set(key, "pending", 0)
		jobId++
		return err
	})
	t.RunEvery(time.Second * 1)
	a.registerTask(t)
	return t
}

func (a *Agent) watchForJobOffersT() *task.Task {
	t := task.New("jobOffersWatcher", func() error {
		key := fmt.Sprintf("%s/%s", JOBS_DIR, "agent_1")

		a.jobC = make(chan *etcd.Response, 1)
		a.jobStopC = make(chan bool, 1)
		a.etcdClient.Watch(key, 0, true, a.jobC, a.jobStopC)
		for {
			r, ok := <-a.jobC
			if !ok {
				a.log("jobC has been closed")
				break
			}
			// Job Received
			a.logf("Job received: %s", r.Node.Key)
		}
		return nil
	})
	t.OnStopFn(func() {
		a.jobStopC <- true
	})
	t.RunOnce()
	a.registerTask(t)
	return t
}

func (a *Agent) watchForNewLeaderElectionT() *task.Task {
	receiverC = make(chan *etcd.Response, 1)
	watchLeaderStopC = make(chan bool, 1)
	t := task.New("watchLeaderElection", func() error {
		a.etcdClient.Watch(SCHEDULER_ELECTION_KEY, 0, true, receiverC, stop)
		<-receiverC
		close(t.ErrorChan()) // future reads of the chan will return nil
	})

	t.OnStopFn(func() {
		watchLeaderStopC <- true
	})
	t.RunOnce()
	a.registerTask(t)
	return t
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

func (a *Agent) registerTask(newTask *task.Task) {
	if len(a.taskManager) == 0 {
		// internal to the task manager that captures cancelation signals
		// from users and trigger the stop of all tasks
		t := task.New("uiTask", func() error {
			<-a.haltTaskC
			return task.ErrUserCanceled
		})
		t.RunOnce()
		a.taskManager = append(a.taskManager, t)
	}
	if a.verbose {
		a.logf("New Task registered: %s", newTask.Id)
	}
	a.taskManager = append(a.taskManager, newTask)
}

func (a *Agent) ListenUserCancelTask() *task.Task {
	return a.taskManager[0]
}

func (a *Agent) Stop() {
	a.haltTaskC <- struct{}{}
	a.haltedC = make(chan struct{}, 1)
	<-a.haltedC
	a.logf("Agent halted")
}

func (a *Agent) stopTasks() {
	a.logf("Stopping all running tasks...")

	var done sync.WaitGroup
	done.Add(len(a.taskManager)) //_taskManager does not need to stop

	for _, t := range a.taskManager {
		go func(task *task.Task) {
			a.logf("Task: %s stopping...", task.Id)
			task.Stop()
			a.logf("Task: %s stopped", task.Id)
			done.Done()
		}(t)
	}

	a.logf("Waiting for %d tasks to stop...", len(a.taskManager))
	done.Wait()

	a.taskManager = a.taskManager[:0]
	a.logf("Sending signal to halt...")
	a.haltedC <- struct{}{}

}
