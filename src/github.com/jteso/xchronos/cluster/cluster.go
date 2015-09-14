// This package contains all the functionality related with the inter-agent communication
// across containers on different hosts, or same host.
package cluster

import (
	"github.com/jteso/xchronos/scheduler"
)

type ClusterClient interface {
	// Connect to the discovery service
	Connect() error

	// Disconnect from the discovery service
	Disconnect()

	// SchedulerElect function attempts to put a lock into the scheduler leadership key. It returns:
	// (true, nil) iff invoker will become the leader
	// (false, nil) iff invoker will become a supporter
	SchedulerElect(ip string) (bool, error)

	// It notifies when the scheduler election key has expired, ie. agent with scheduler role has failed
	// to renew the key
	SchedulerFailureWatcher(notify chan bool, stopC chan bool)

	// Publish a job due to be executed immediately
	MakeJobOffer(*scheduler.Job) error

	// It notifes when a new offer has been published
	WatchJobOffers(notify chan *scheduler.Job, stopC chan bool)

	// Declare the agent's intention to execute the job, by registering the agents ip
	// return true if successful
	TakeJobOffer(job *scheduler.Job, ip string) bool

	// Register the agent as an executor
	RegisterAsExecutor(agentId, agentIp string) error

	// Persist all client scheduled jobs
	PersistJob(*scheduler.Job) error
}
