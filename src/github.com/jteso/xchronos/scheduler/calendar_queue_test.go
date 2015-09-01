package scheduler

import (
	"testing"
	"time"
)

func now() time.Time {
	return time.Now()
}

func TestBucketAllocation(t *testing.T) {
	bucketWidth := int64(3600) // 1 hour
	bucketLen := 3
	calq := NewCalq(bucketWidth, bucketLen)

	// job due now
	now := time.Now()
	assert(t, calq.allocateBucket(now.UnixNano()) == 0, "it should be zero")
	assert(t, calq.allocateBucket(now.Add(61*time.Minute).UnixNano()) == 1, "it should be one")
	assert(t, calq.allocateBucket(now.Add(121*time.Minute).UnixNano()) == 2, "it should be two")
	assert(t, calq.allocateBucket(now.Add(181*time.Minute).UnixNano()) == 0, "it should be zero")
	assert(t, calq.allocateBucket(now.Add(241*time.Minute).UnixNano()) == 1, "it should be one")
}

func TestCalendarqueue(t *testing.T) {
	bucketWidth := int64(3600) // 1 hour
	bucketLen := 3
	calq := NewCalq(bucketWidth, bucketLen)

	now := time.Now()
	// jobs
	job1 := &Job{
		Id:        "1",
		nextRunAt: now.Add(-2 * time.Second),
	}
	job2 := &Job{
		Id:        "2",
		nextRunAt: now.Add(61 * time.Minute),
	}

	calq.Enqueue(job1)
	calq.Enqueue(job2)

	j := calq.Next()

	assert(t, j != nil, "Job cannot be nil")
	assert(t, j.Id == "1", "Job.id expected <%s> actual <%s>", job1.Id, j.Id)

	// keep in mind that next is just peeking the value, not popping it
	j = calq.Next()
	assert(t, j.Id == "1", "Job.id expected <%s> actual <%s>", job2.Id, j.Id)

	calq.Dequeue(j)

	//j2 := calq.Next()
	//assert(t, j2.Id == "2", "Job.id expected <%s> actual <%s>", job2.Id, j2.Id)
}

func TestSkipFutureJobsCurrentBucket(t *testing.T) {
	bucketWidth := int64(3600) // 1 hour
	bucketLen := 3
	calq := NewCalq(bucketWidth, bucketLen)

	// bucket 0
	job1 := &Job{
		Id:        "1",
		nextRunAt: now().Add(-2 * time.Second),
	}
	// bucket 1
	job2 := &Job{
		Id:        "2",
		nextRunAt: now().Add(61 * time.Minute),
	}
	// bucket 0
	job3 := &Job{
		Id:        "3",
		nextRunAt: now().Add(181 * time.Minute),
	}

	calq.Enqueue(job1)
	calq.Enqueue(job2)
	calq.Enqueue(job3)

	// Execution order should be job1, job2, job3
	// So ensure the bucket 0 will not schedule the job3 before the job2
	j1 := calq.Next()
	assert(t, j1.Id == job1.Id, "Job.id expected <%s> actual <%s>", job1.Id, j1.Id)
	calq.Dequeue(j1)

	j2 := calq.Next()
	assert(t, j2.Id == job2.Id, "Job.id expected <%s> actual <%s>", job2.Id, j2.Id)
	calq.Dequeue(j2)

	j3 := calq.Next()
	assert(t, j3.Id == job3.Id, "Job.id expected <%s> actual <%s>", job3.Id, j3.Id)
	calq.Dequeue(j3)
}
