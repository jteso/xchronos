// This package contains all the functionality related with the inter-agent communication
// across containers on different hosts, or same host.
package cluster

type ClusterClient interface {
	// Connect to the discovery service
	Connect() error

	// Disconnect from the discovery service
	Disconnect()

	// SchedulerElect function attempts to put a lock into the scheduler leadership key. It returns:
	// if successful(true, nil): invoker will become the leader
	// otherwise, if (false, nil): invoker will become a supporter
	SchedulerElect(ip string) (bool, error)

	NotifySchedulerStepDown(notify chan bool, stopC chan bool)

	// Publish a job due to be executed immediately
	MakeJobOffer(*Job) error

	// Get notification upon publication of new jobs
	WatchJobOffers(notify chan *Job, stopC chan bool)

	// Declare the agent's intention to execute the job
	TakeJobOffer() (*Job, error)

	// Register the agent as an executor
	RegisterAsExecutor(*Agent) error

	// Persist all client scheduled jobs
	PersistJob(*Job) error
}
