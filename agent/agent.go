package agent

import (
	"fmt"
	"log"
	"time"
	"github.com/jteso/xchronos/Godeps/_workspace/src/github.com/coreos/go-etcd/etcd"
	"strconv"
)

const (
	SCHEDULER_ELECTION_KEY = "/xchronos/var/scheduler/election"
	EXECUTORS_DIR = "/xchronos/etc/executors"
	JOBS_DIR = "/xchronos/etc/jobs"
)

var(
	SCHEDULER_LEADER_TTL uint64  = 10 // max time running without leader
	EXECUTOR_TTL uint64 = 10		  // max time running without a particular executor
	HEARTBEAT time.Duration = 5
)


type Agent struct {
	ID         string
	state     string
	etcdClient *etcd.Client
	etcdNodes  []string
	verbose bool

	executorTicker *time.Ticker
}

// etcdNotes = strings.Split(os.Getenv("ETCD_NOTES"), ",")
func New(id string, etcdNodes []string, verbose bool) *Agent {
	return &Agent{
		ID:        id,
		state: 	   "INIT",
		etcdNodes: etcdNodes,
		verbose: verbose,
	}
}

// panic function will force an agent to stop its execution, due to some
// unrecoverable error found.
func (a *Agent) panic(format string, args ...interface{}) handleStateFn {
	fmt.Printf(format, args)
	return nil
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

func (a *Agent) changeState(newState string){
	a.log(fmt.Sprintf("Changing state: %s -> %s", a.state, newState))
	a.state = newState
}

// registerAsExecutor function will make the agent to offer itself to execute jobs been offered.
// this offering is been done by writing periodically (HEARTBEAT) into the etcd dir /executors
// this function will returned via chan any error is encountered, and this agent can be stop been
// offered as an executor by stopping the agents executorTicker.
func (a *Agent) registerAsExecutor() (errCh chan error){
	errCh = make(chan error, 1)
	key := EXECUTORS_DIR + "/" + a.ID

	a.log("Registered as executor. key=" + key)

	var err error
	a.executorTicker = time.NewTicker(time.Second * HEARTBEAT)
	go func() {
		for {
			<- a.executorTicker.C
			_, err = a.etcdClient.Set(key, "up" , EXECUTOR_TTL)
			if err != nil {
				a.executorTicker.Stop()
				errCh <- err
			}
		}
	}()

	return errCh
}

func (a *Agent) publishJobOffers() (errCh chan error){
	errCh = make(chan error, 1)
	jobId := 0
	var key string

	a.log("Waiting for scheduler to start publishing jobs...")

	var err error
	jobTicker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			<- jobTicker.C
			key = JOBS_DIR + "/agent_1/" + strconv.Itoa(jobId)
			a.log("Publishing job: "+ key)
			_, err = a.etcdClient.Set(key, "pending", 0)
			jobId ++

			if err != nil {
				jobTicker.Stop()
				errCh <- err
			}
		}
	}()

	return errCh
}

func (a *Agent) watchForJobOffers() (errCh chan error){
	watchChan := make(chan *etcd.Response)
	key := fmt.Sprintf("%s/%s", JOBS_DIR, "agent_1")
	a.log(fmt.Sprintf("Waiting for jobs on: %s", key))
	go a.etcdClient.Watch(key, 0, true, watchChan, nil)
	for {
		r := <-watchChan
		a.log(fmt.Sprintf("Job received: %s", r.Node.Key))
	}

	return nil
}

func (a *Agent) log(message string) {
	if a.verbose {
		log.Printf("[%s] %s\n", a.ID, message)

	}
}


