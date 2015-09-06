package scheduler

import (
	"math/rand"
	"strconv"
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

// -------
// Benchmark for calendar queue
// -------

// go test -bench=. (it will execute tests)
// got test -run=XXX -bench=.  -benchtime=20s (it will skip the tests)
func enqueueRandomEvent(num int) {
	bucketWidth := int64(3600)
	bucketLen := 3
	calq := NewCalq(bucketWidth, bucketLen)

	var job *Job
	for i := 0; i < num; i++ {
		job = &Job{
			Id:        strconv.Itoa(i),
			nextRunAt: now().Add(time.Duration(rand.Intn(100000000)) * time.Minute),
		}
		calq.Enqueue(job)
	}
}

func benchmarkCalq(i int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		enqueueRandomEvent(i) // enqueue i random events
	}
}

// func BenchmarkCalq10(b *testing.B) {
// 	benchmarkCalq(10, b)
// }

// func BenchmarkCalq20(b *testing.B) {
// 	benchmarkCalq(20, b)
// }

// func BenchmarkCalq50(b *testing.B) {
// 	benchmarkCalq(50, b)
// }

// func BenchmarkCalq100(b *testing.B) {
// 	benchmarkCalq(100, b)
// }

// func BenchmarkCalq500(b *testing.B) {
// 	benchmarkCalq(500, b)
// }

// func BenchmarkCalq1000(b *testing.B) {
// 	benchmarkCalq(1000, b)
// }
// func BenchmarkCalq2000(b *testing.B) {
// 	benchmarkCalq(2000, b)
// }

func BenchmarkCalq5000(b *testing.B) {
	benchmarkCalq(5000, b)
}

func BenchmarkCalq50000(b *testing.B) {
	benchmarkCalq(50000, b)
}

func BenchmarkCalq500000(b *testing.B) {
	benchmarkCalq(500000, b)
}

func BenchmarkCalq1000000(b *testing.B) {
	benchmarkCalq(1000000, b)
}

func BenchmarkCalq100000000(b *testing.B) {
	benchmarkCalq(100000000, b)
}
