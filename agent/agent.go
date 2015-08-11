package agent

import(
	"fmt"
	"time"
)

type Operation interface{
	GetAckCh() chan bool
	WaitForAck()
}

type TerminateOp struct {
	ack chan bool
}
func (this TerminateOp) GetAckCh() chan bool {
	return this.ack
}
func (this TerminateOp) WaitForAck() {
	<- this.ack
}

func NewTerminateOp() Operation {
	return &TerminateOp{
		ack: make(chan bool),
	}
}

type AckableOpCh chan Operation


type AgentState int

const(
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
	Name string
	stopCh AckableOpCh
}
func New(name string, stopCh AckableOpCh) *Agent {
	return &Agent{
		Name: name,
		stopCh: stopCh,
	}
}

func (a *Agent) errorf(format string, args ...interface{}) handleStateFn {
	fmt.Printf(format, args)
	return nil
}

func (a *Agent) Run() error {
	for stateHandler := initAgent; stateHandler != nil; {
		stateHandler = stateHandler(a)
	}
	close(a.stopCh)
	return nil
}

func initAgent(agent *Agent) handleStateFn {
	return joinEtcdCluster
}

func joinEtcdCluster(agent *Agent) handleStateFn {
	fmt.Printf("Agent: %s joining etcd cluster...\n", agent.Name)
	time.Sleep(2 * time.Second)
	return electForLeader
}

func electForLeader(agent *Agent) handleStateFn {
	fmt.Printf("Agent: %s electing for leader\n", agent.Name)
	time.Sleep(2 * time.Second)
	return nil
}



