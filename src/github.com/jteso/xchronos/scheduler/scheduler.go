package scheduler

import (
	"log"
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

func (s *Scheduler) Notify(dueJobC chan *Job, stopC chan bool) {
	go func() {
		for {
			select {
			case <-stopC:
				log.Printf("Scheduler received stop notification\n")
				close(dueJobC)
				return
			default:
				job := s.dequeue()
				if job != nil {
					//log.Printf("Next job due for execution in %f secs...\n", job.WaitSecs())
					time.Sleep(time.Duration(job.WaitSecs()) * time.Second)
					dueJobC <- job
				} else {
					//log.Printf("No more jobs available\n")
				}
			}
		}
	}()
}

func (s *Scheduler) dequeue() *Job {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	job, _ := s.jobQueue.Pop()
	return job
}

// TODO(javier): check whether this job has been scheduled already before
func (s *Scheduler) Enqueue(job *Job) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.jobQueue.Push(job)
}

func (s *Scheduler) Size() int {
	return s.jobQueue.Size()
}
