package cluster

import (
	"fmt"
	"time"

	"github.com/coreos/go-etcd/etcd"
	"github.com/jteso/xchronos/scheduler"
)

const (
	SCHEDULER_ELECTION_KEY    = "/zeitd/lock/scheduler"
	JOB_OFFER_KEY             = "/zeitd/offers"
	JOB_TAKER_KEY             = "/zeitd/lock/job"
	EXECUTOR_REGISTRATION_KEY = "/zeitd/executor/ip"
	JOB_PERSISTENCE_KEY       = "/zeitd/jobs/job"
)

var (
	SCHEDULER_LEADER_TTL uint64        = 10 // max time running without leader
	EXECUTOR_TTL         uint64        = 10 // max time running without a particular executor
	HEARTBEAT            time.Duration = 5
	NO_TTL                             = uint64(0)
)

type EtcdProxy struct {
	// Used to communicate to etcd cluster
	etcdClient *etcd.Client
	// See all the etcd nodes available in the cluster
	etcdNodes []string
}

func NewEtcdClient(etcdNodes []string) *EtcdProxy {
	proxy := new(EtcdProxy)
	proxy.etcdNodes = etcdNodes
	proxy.etcdClient = etcd.NewClient(etcdNodes)

	return proxy
}

func (e *EtcdProxy) Connect() error {
	// TODO: some sort of ping/version/...
	return nil
}

func (e *EtcdProxy) Disconnect() {
	e.etcdClient.Close()
}

func (e *EtcdProxy) SchedulerElect(ip string) (bool, error) {
	_, err := e.etcdClient.Create(SCHEDULER_ELECTION_KEY, ip, SCHEDULER_LEADER_TTL)
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

func (e *EtcdProxy) SchedulerFailureWatcher(notifyC chan bool, stopC chan bool) {
	receiverC := make(chan *etcd.Response, 1)
	e.etcdClient.Watch(SCHEDULER_ELECTION_KEY, 0, true, receiverC, stopC)
	<-receiverC
	notifyC <- true
}

// Make a job offer in: (k,v) -> ($JOB_OFFERS/$job_id , job)
func (e *EtcdProxy) MakeJobOffer(job *scheduler.Job) error {
	key := fmt.Sprintf("%s/%s", JOB_OFFER_KEY, job.Id)
	jobAsBytes, _ := job.Bytes()
	_, err := e.etcdClient.Set(key, string(jobAsBytes), NO_TTL)
	return err
}

// TODO: Reimplementation required
func (e *EtcdProxy) WatchJobOffers(notify chan *scheduler.Job, stopC chan bool) {
	// jobC := make(chan *etcd.Response, 1)
	// defer close(jobC)
	// e.etcdClient.Watch(JOB_OFFER_KEY, 0, true, jobC, stopC)
	// for {
	// 	r, ok := <-jobC
	// 	if !ok {
	// 		log.Print("jobC has been closed")
	// 		break
	// 	}
	// 	// Job Received
	// 	log.Printf("Offer registered: %s", r.Node.Key)
	// 	job := NewFromBytes([]byte(r.Node.Value))
	// 	notify <- job
	// }

}

func (e *EtcdProxy) TakeJobOffer(job *scheduler.Job, ip string) bool {
	if e.acquireLockToJobOffer(job.Id, ip) {
		// remove the job offer
		e.etcdClient.Delete(fmt.Sprintf("%s/%s", JOB_OFFER_KEY, job.Id), false)
		// remove the lock
		e.etcdClient.Delete(fmt.Sprintf("%s/%s", JOB_TAKER_KEY, job.Id), false)
		return true
	}
	return false
}

func (e *EtcdProxy) RegisterAsExecutor(agentId, agentIp string) error {
	_, err := e.etcdClient.Set(fmt.Sprintf("%s/%s", EXECUTOR_REGISTRATION_KEY, agentId), agentIp, EXECUTOR_TTL)
	return err
}

func (e *EtcdProxy) PersistJob(*scheduler.Job) error {
	//TODO
	return nil
}

func (e *EtcdProxy) acquireLockToJobOffer(jobId string, takerIp string) bool {
	etcdLockPath := fmt.Sprintf("%s/%s", JOB_TAKER_KEY, jobId)
	_, err := e.etcdClient.Create(etcdLockPath, takerIp, NO_TTL)
	if err != nil {
		if etcdError, ok := err.(*etcd.EtcdError); ok {
			// key already exists
			if etcdError.ErrorCode == 105 {
				return false
			}
		}
		return false
	}
	return true
}
