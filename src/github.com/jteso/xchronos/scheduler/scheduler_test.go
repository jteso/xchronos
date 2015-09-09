package scheduler

import (
	"fmt"
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

	// Dequeue all the mock jobs
	var job *Job
	for j := 0; j < totalJobs; j++ {
		job = sch.NextJob()
		fmt.Printf("--> Received job_id=%s\n", job.Id)
		assert.Equal(t, job.Id, strconv.Itoa(j), "Received the wrong job")
	}

	// Make sure no more jobs are pending
	next := sch.NextJob()
	assert.Nil(t, next)
}

func jobGenerator(n int) []*Job {
	output := make([]*Job, n)
	for i := 0; i < n; i++ {
		output[i] = NewTestJob(strconv.Itoa(i), i+5)
	}
	return output
}
