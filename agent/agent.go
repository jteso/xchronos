package agent

import (
	"fmt"
	"github.com/jteso/xchronos/Godeps/_workspace/src/github.com/coreos/go-etcd/etcd"
	"time"
)

const (
	SCHEDULER_ELECTION_KEY = "/xchronos/scheduler/election"
	SCHEDULER_LEADER_TTL   = 5
)

type AgentState int

const (
	STATE_INIT AgentState = iota
	STATE_ELECT_FOR_LEADER
	STATE_NODE_BECOME_LEADER
	STATE_NODE_BECOME_FOLLOWER
	STATE_NODE_ERROR
)

// agentStateFn represents the state of an agent as a function
// that returns the next state
type handleStateFn func(*Agent) handleStateFn

type Agent struct {
	ID         string
	etcdClient *etcd.Client
	etcdNodes  []string
}

// etcdNotes = strings.Split(os.Getenv("ETCD_NOTES"), ",")
func New(id string, etcdNodes []string) *Agent {
	return &Agent{
		ID:        id,
		etcdNodes: etcdNodes,
	}
}

// panic function will force an agent to stop its execution, due to some
// unrecoverable error found.
func (a *Agent) panic(format string, args ...interface{}) handleStateFn {
	fmt.Printf(format, args)
	return nil
}

func (a *Agent) Run() error {
	for stateHandler := initAgent; stateHandler != nil; {
		stateHandler = stateHandler(a)
	}
	//close(a.stopCh)
	return nil
}

func initAgent(agent *Agent) handleStateFn {
	return joinEtcdCluster
}

func joinEtcdCluster(agent *Agent) handleStateFn {
	fmt.Printf("[Agent:%s] Joining etcd cluster...\n", agent.ID)
	agent.etcdClient = etcd.NewClient(agent.etcdNodes)

	time.Sleep(2 * time.Second)
	return runningForLeader
}

// TODO - use the Atomic Compare and Swap (CAS) for the locking service
func runningForLeader(agent *Agent) handleStateFn {
	fmt.Printf("[Agent:%s] Running for leader...\n", agent.ID)

	// Put value if prevExist=false
	_, err := agent.etcdClient.Create(SCHEDULER_ELECTION_KEY, agent.ID, SCHEDULER_LEADER_TTL)

	if err != nil {
		if etcdError, ok := err.(*etcd.EtcdError); ok {
			// key already exists
			if etcdError.ErrorCode == 105 {
				return actAsSupporter
			}else{
				return agent.panic("Etcd returned unexpected error: %s", etcdError.Error())
			}
		}else{
			return agent.panic("Error while electing for leadership due to: %s", err.Error())
		}
		
	}

	return actAsLeader
}

func actAsLeader(agent *Agent) handleStateFn {
	fmt.Printf("[Agent:%s] stepping up as leader\n", agent.ID)
	// ensure periodic leadership renewal
	// setup crontab in mem on local node
	// [cron signal] publish job offers
	// [etcd leader watcher] return electForLeader
	// [etcd jobOffer watcher] competeForJobExecution
	return nil
}

func actAsSupporter(agent *Agent) handleStateFn {
	fmt.Printf("[Agent:%s] stepping down\n", agent.ID)
	// [etcd leader watcher] return electForLeader
	// [etcd jobOffer watcher] competeForJobExecution
	return nil
}
