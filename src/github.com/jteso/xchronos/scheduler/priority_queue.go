// This file is a modified version over the original
// Source: https://raw.githubusercontent.com/oleiade/lane/master/pqueue.go
// Author: oleiade (github)
package scheduler

import (
	"fmt"
	"sync"
)

// PQType represents a priority queue ordering kind (see MAXPQ and MINPQ)
type PQType int

const (
	MAXPQ PQType = iota
	MINPQ
)

type item struct {
	// Job
	job *Job
	// number of secs elapsed since Jan 1, 1970 UTC
	priority int64
}

// PQueue is a heap priority queue data structure implementation.
// It can be whether max or min ordered and it is synchronized
// and is safe for concurrent operations.
type PQueue struct {
	sync.RWMutex
	items      []*item
	elemsCount int
	comparator func(int64, int64) bool
}

func newItem(job *Job, priority int64) *item {
	return &item{
		job:      job,
		priority: priority,
	}
}

// NewPQueue creates a new priority queue with the provided pqtype
// ordering type
func NewPQueue(pqType PQType) *PQueue {
	var cmp func(int64, int64) bool

	if pqType == MAXPQ {
		cmp = max
	} else {
		cmp = min
	}

	items := make([]*item, 1)
	items[0] = nil // Heap queue first element should always be nil

	return &PQueue{
		items:      items,
		elemsCount: 0,
		comparator: cmp,
	}
}

// Push the job item into the priority queue with provided priority.
func (pq *PQueue) Push(job *Job, priority int64) {
	item := newItem(job, priority)

	pq.Lock()
	pq.items = append(pq.items, item)
	pq.elemsCount += 1
	pq.swim(pq.size())
	pq.Unlock()
}

// Pop and returns the highest/lowest priority item (depending on whether
// you're using a MINPQ or MAXPQ) from the priority queue
func (pq *PQueue) Pop() (*Job, int64) {
	pq.Lock()
	defer pq.Unlock()

	if pq.size() < 1 {
		return nil, 0
	}

	var max *item = pq.items[1]

	pq.exch(1, pq.size())
	pq.items = pq.items[0:pq.size()]
	pq.elemsCount -= 1
	pq.sink(1)

	return max.job, max.priority
}

// Head returns the highest/lowest priority item (depending on whether
// you're using a MINPQ or MAXPQ) from the priority queue
func (pq *PQueue) Head() (*Job, int64) {
	pq.RLock()
	defer pq.RUnlock()

	if pq.size() < 1 {
		return nil, 0
	}

	headValue := pq.items[1].job
	headPriority := pq.items[1].priority

	return headValue, headPriority
}

// Size returns the elements present in the priority queue count
func (pq *PQueue) Size() int {
	pq.RLock()
	defer pq.RUnlock()
	return pq.size()
}

func (pq *PQueue) size() int {
	return pq.elemsCount
}

func max(i, j int64) bool {
	return i < j
}

func min(i, j int64) bool {
	return i > j
}

func (pq *PQueue) less(i, j int) bool {
	return pq.comparator(pq.items[i].priority, pq.items[j].priority)
}

func (pq *PQueue) exch(i, j int) {
	var tmpItem *item = pq.items[i]

	pq.items[i] = pq.items[j]
	pq.items[j] = tmpItem
}

func (pq *PQueue) swim(k int) {
	for k > 1 && pq.less(k/2, k) {
		pq.exch(k/2, k)
		k = k / 2
	}

}

func (pq *PQueue) sink(k int) {
	for 2*k <= pq.size() {
		var j int = 2 * k

		if j < pq.size() && pq.less(j, j+1) {
			j++
		}

		if !pq.less(k, j) {
			break
		}

		pq.exch(k, j)
		k = j
	}
}

func (pq *PQueue) String() string {
	return fmt.Sprintf("[%+v]", pq.items)
}
