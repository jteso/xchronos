package scheduler

import "time"

//github.com/gorhill/cronexpr
type JobPublisher interface {
	Publish(*Job) error
}

type EtcdPublisher struct {
	etcdPubPath string
}

func NewEtcdPublisher(pubPath string) *EtcdPublisher {
	return &EtcdPublisher{
		etcdPubPath: pubPath,
	}
}
func (this *EtcdPublisher) Publish(j *Job) error {
	// TODO: publish job onto etcd
}

type Scheduler struct {
	// Cache of all jobs
	jobQueue *PQueue
	// Chan to receive new jobs to schedule
	newJobsC chan *Job
	// job publisher
	jobPublisher JobPublisher
}

func New() *Scheduler {
	return &Scheduler{
		jobQueue:     NewPQueue(MAXPQ),
		newJobsC:     make(chan *Job, 10),
		jobPublisher: NewEtcdPublisher(pubPath),
	}
}

func (s *Scheduler) Start() error {
	go s.jobReceive()
	go s.jobPublisher()
}

func (s *Scheduler) EnqueueJob(j *Job) {
	s.newJobsC <- j
}

func (s *Scheduler) jobReceive(j *Job) {
	for {
		select {
		case j <- s.newJobsC:
			s.jobQueue.Push(j, j.GetNextRunAt())
		}
	}
}

func (s *Scheduler) jobPublisher() {
	for {
		job, _ := s.jobQueue.Pop()
		waitSec := job.WaitSecs()
		time.Sleep(waitSec)
		s.jobPublisher.Publish(job)
	}
}
