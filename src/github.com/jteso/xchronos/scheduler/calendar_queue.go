// This file contains an implementation of a calendar queue.
// Specification: http://pioneer.netserv.chula.ac.th/~achaodit/paper5.pdf
package scheduler

type bucket *PQueue

type Calq struct {
	// general counter
	totalEvents int64
	// Admission time range for each bucket
	bucketWidth int64
	// List of priority queues, a bucket is defined as  a list with a
	// specified range of admission times
	buckets []*bucket
}

// Optimal values: M=2N (num buckets equals double of number of events)
// N=100 events
// M=200
// w=432 (sec)
// hence,
// 1gen buckets = todays
// 2gen buckets = tomorrow
// igen buckets = today + (i-1)
func NewCalq() *Calq {
	return &Calq{
		totalEvents: 0,
		bucketWidth: 60 * 60, //1hour
		buckets:     make([]*bucket, 24),
	}
}

// Find the next event in the calendar queue to execute
func (c *Calq) Next() *Job {
	return nil
}

// Enqueue a new event, or recently executed one,into the calendar queue
// Algorithm:
// To find the bucket number m(e) to enqueue an event e that occurs at time t(e), we may apply
// m(e) = [t(e) / bucketWidth] mod M

func (c *Calq) Enqueue(j *Job) {

}
