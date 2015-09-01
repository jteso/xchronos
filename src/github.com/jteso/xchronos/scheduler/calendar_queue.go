// This file contains an implementation of a calendar queue.
// Specification: http://pioneer.netserv.chula.ac.th/~achaodit/paper5.pdf
package scheduler

import (
	"fmt"
	"sync"
	"time"
)

// A bucket is defined as  a list with a
// specified range of admission times
type Bucket interface {
	Push(job *Job, priority int64)
	Pop() (*Job, int64)
	Head() (*Job, int64)
	String() string
}

// -- Calendar Queue definition
type Calq struct {
	// general counter
	totalEvents int64
	mutexEvents *sync.RWMutex

	// Admission time range for each bucket in seconds
	bucketWidth int64
	// List of buckets,
	buckets []Bucket
}

// Optimal values: bucketLen = 2 * number_events_to_manage
func NewCalq(bucketWidth int64, bucketLen int) *Calq {
	// initialize the buckets with their corresponding priority queues
	bkts := make([]Bucket, bucketLen)
	for i, _ := range bkts {
		bkts[i] = NewPQueue(MINPQ)
	}

	return &Calq{
		totalEvents: 0,
		mutexEvents: &sync.RWMutex{},
		bucketWidth: bucketWidth, //60 * 60 (1 hour)
		buckets:     bkts,
	}
}

// For debugging purposes only
func (c *Calq) Print() {
	for i, pq := range c.buckets {
		fmt.Printf("[bucket: %d]\t%s\n", i, pq.String())
	}

}

// Find the next event in the calendar queue to execute
func (c *Calq) Next() *Job {
	// Trival case
	c.mutexEvents.RLock()
	if c.totalEvents == 0 {
		return nil
	}
	c.mutexEvents.Unlock()

	// lets iterate over all the buckets of the calendar queue, starting from current bucket (relative to time.Now) and from cycle 0 onwards
	return c.findNextJob(0, c.headBucket(), 0)
}

// Enqueue a new event, or recently executed one,into the calendar queue
// Algorithm:
// To find the bucket number m(e) to enqueue an event e that occurs at time t(e), we may apply
// m(e) = [t(e) / bucketWidth] mod M

func (c *Calq) Enqueue(job *Job) {
	// Find out the next run date in nanoseconds
	scheduledNs := job.GetNextRunAt().UnixNano()
	// Resolve the bucket
	bucketId := c.allocateBucket(scheduledNs)
	// Enqueue it in the priority queue for that bucket
	c.buckets[bucketId].Push(job, job.GetNextRunAt().UnixNano())
	// Increasing number of events
	c.mutexEvents.Lock()
	c.totalEvents = c.totalEvents + 1
	c.mutexEvents.Unlock()

}

// Just remove the job
func (c *Calq) Dequeue(job *Job) {
	// Find out the next run date in nanoseconds
	scheduledNs := job.GetNextRunAt().UnixNano()
	// Resolve the bucket
	bucketId := c.allocateBucket(scheduledNs)
	// Enqueue it in the priority queue for that bucket
	c.buckets[bucketId].Pop()
	// Decreasing number of events
	c.mutexEvents.Lock()
	c.totalEvents = c.totalEvents - 1
	c.mutexEvents.Unlock()
}

//return the bucket where the next job is to be run
func (c *Calq) headBucket() int {
	return c.allocateBucket(time.Now().UnixNano())
}

// Find out the bucket where a particular execution date (unixNano) belongs to.
func (c *Calq) allocateBucket(unixNano int64) int {
	// Calculate when is due to execute in secs
	runOnSec := (unixNano - time.Now().UnixNano()) / time.Second.Nanoseconds()
	//fmt.Printf("-- runOnSec: %d\n", runOnSec)
	x := (runOnSec / c.bucketWidth) % int64(len(c.buckets))

	return int(x)
}

// maxTimeCurrentBucketNs returns in nanoseconds the max time a job should be scheduled
// in order to be considered to run on the current cycle.
// This function is needed as a bucket include jobs within a particular time range, but
// if the job is scheduled for a far away in the future, it may end up in the same bucket.
// If that occurs we should skip the current bucket and start processing the next one.
//
// Example:
// [bucket i] -> (job(5 min), job(10 min), job(1 year))
// [bucket j] -> (job(30 min))
//
// In the example above, you can notice how the job(1 year) need to be skipped as the the job
// does not belong to current execution cycle.
func (c *Calq) belongsCurrentCycle(job *Job, cycle int) bool {
	bucketWidhSec := time.Duration(c.bucketWidth) * time.Second
	calqWidthSec := time.Duration(c.bucketWidth*int64(len(c.buckets)*cycle)) * time.Second
	upperLimit := (time.Now().UnixNano() + bucketWidhSec.Nanoseconds()) + calqWidthSec.Nanoseconds()
	return job.GetNextRunAt().UnixNano() < upperLimit
}

func (c *Calq) findNextJob(cycle int, currentBkt int, scannedBkts int) *Job {
	jobNext, _ := c.buckets[currentBkt].Head()
	for jobNext == nil {
		currentBkt = (currentBkt + 1) % len(c.buckets) // bkt++ with upper-bound control
		jobNext, _ = c.buckets[currentBkt].Head()
	}
	// <-- start
	if scannedBkts == len(c.buckets) {
		return c.findNextJob(cycle+1, currentBkt, 0)
	} else {
		jobNext, _ := c.buckets[currentBkt].Head()
		// bucket empty or does not belong to current cycle, then iterate over the next bucket
		if jobNext == nil || !c.belongsCurrentCycle(jobNext, cycle) {
			return c.findNextJob(cycle, (currentBkt+1)%len(c.buckets), scannedBkts+1)
		} else {
			return jobNext
		}
	}

}
