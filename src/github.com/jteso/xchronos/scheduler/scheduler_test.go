package scheduler

import (
	"strconv"
	"testing"

	"github.com/jteso/testify/assert"
)

func TestHappyScheduler(t *testing.T) {
	var totalJobs = 3

	// Start the scheduler
	sch := NewScheduler()

	// Enqueue all the mock jobs
	mockJobs := jobGenerator(totalJobs)
	for i := 0; i < totalJobs; i++ {
		sch.Enqueue(mockJobs[i])
	}

	assert.Equal(t, sch.Size(), 3)

	// Dequeue all the mock jobs
	cont := 0
	dueJobC := make(chan *Job, 10)
	stopC := make(chan bool)
	sch.Notify(dueJobC, stopC)

	for job := range dueJobC {
		assert.Equal(t, job.Id, strconv.Itoa(cont), "Received the wrong job")
		cont += 1
		if cont == totalJobs {
			stopC <- true
		}
	}

	assert.Equal(t, sch.Size(), 0)
}

func jobGenerator(n int) []*Job {
	output := make([]*Job, n)
	for i := 0; i < n; i++ {
		output[i] = NewTestJob(strconv.Itoa(i), i+5)
	}
	return output
}
