package scheduler

import (
	"sync"
	"time"
)

type Scheduler struct {
	// Cache of all jobs
	jobQueue *PQueue
	mutex    *sync.Mutex
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		jobQueue: NewPQueue(MINPQ),
		mutex:    &sync.Mutex{},
	}
}

func (s *Scheduler) NextJob() *Job {
	return s.popAndWait()
}

func (s *Scheduler) Enqueue(job *Job) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.jobQueue.Push(job)
}

func (s *Scheduler) popAndWait() *Job {
	s.mutex.Lock()
	job, _ := s.jobQueue.Pop()
	if job == nil {
		return nil
	}
	s.mutex.Unlock()
	waitSec := job.WaitSecs()
	time.Sleep(time.Duration(waitSec) * time.Second)
	return job

}
