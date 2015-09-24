package scheduler

import (
	_ "math/rand"
	_ "reflect"
	_ "strconv"
	_ "sync"
	_ "testing"
	_ "time"

	_ "github.com/jteso/testify/assert"
)

// func TestMaxPQueue_init(t *testing.T) {
// 	pqueue := NewPQueue(MAXPQ)

// 	assert.Equal(t, len(pqueue.items), 1, "incorrect queue length")
// 	assert.Equal(t, pqueue.Size(), 0, "incorrect queue length")
// 	assert.Nil(t, pqueue.items[0])

// 	assert.Equal(t, reflect.ValueOf(pqueue.comparator).Pointer(), reflect.ValueOf(max).Pointer(), "pqueue.comparator != max")
// }

// func TestMinPQueue_init(t *testing.T) {
// 	pqueue := NewPQueue(MINPQ)

// 	assert.Equal(t, len(pqueue.items), 1)
// 	assert.Equal(t, pqueue.Size(), 0)
// 	assert.Nil(t, pqueue.items[0])
// 	assert.Equal(t, reflect.ValueOf(pqueue.comparator).Pointer(), reflect.ValueOf(min).Pointer())
// }

// func TestMaxPQueuePushAndPop_protects_max_order(t *testing.T) {
// 	pqueue := NewPQueue(MAXPQ)
// 	pqueueSize := 100

// 	// Populate the test priority queue with dummy elements
// 	// in asc ordered.
// 	for i := 0; i < pqueueSize; i++ {
// 		value := NewTestJob(strconv.Itoa(i), i)
// 		pqueue.Push(value)
// 	}

// 	containerIndex := 1 // binary heap are 1 indexed
// 	for i := 99; i >= 0; i-- {
// 		expectedValue := NewTestJob(strconv.Itoa(i), i)
// 		expectedPriority := int64(i)

// 		// Avoiding testing arithmetics headaches by using the pop function directly
// 		value, priority := pqueue.Pop()
// 		assert.Equal(t, value, expectedValue)
// 		assert.Equal(t, priority, expectedPriority)

// 		containerIndex++
// 	}
// }

// func TestMaxPQueuePushAndPop_concurrently_protects_max_order(t *testing.T) {
// 	var wg sync.WaitGroup

// 	pqueue := NewPQueue(MAXPQ)
// 	pqueueSize := 100

// 	// Populate the test priority queue with dummy elements
// 	// in asc ordered.
// 	for i := 0; i < pqueueSize; i++ {
// 		wg.Add(1)

// 		go func(i int) {
// 			defer wg.Done()

// 			value := NewTestJob(strconv.Itoa(i), i)
// 			pqueue.Push(value)
// 		}(i)
// 	}

// 	wg.Wait()

// 	containerIndex := 1 // binary heap are 1 indexed
// 	for i := 99; i >= 0; i-- {
// 		expectedValue := NewTestJob(strconv.Itoa(i), i)
// 		expectedPriority := int64(i)

// 		// Avoiding testing arithmetics headaches by using the pop function directly
// 		value, priority := pqueue.Pop()
// 		assert.Equal(t, value, expectedValue)
// 		assert.Equal(t, priority, expectedPriority)

// 		containerIndex++
// 	}
// }

// func TestMinPQueuePushAndPop_protects_min_order(t *testing.T) {
// 	pqueue := NewPQueue(MINPQ)
// 	pqueueSize := 100

// 	// Populate the test priority queue with dummy elements
// 	// in asc ordered.
// 	for i := 0; i < pqueueSize; i++ {
// 		value := NewTestJob(strconv.Itoa(i), i)
// 		pqueue.Push(value)
// 	}

// 	for i := 0; i < pqueueSize; i++ {
// 		expectedValue := NewTestJob(strconv.Itoa(i), i)
// 		expectedPriority := int64(i)

// 		// Avoiding testing arithmetics headaches by using the pop function directly
// 		value, priority := pqueue.Pop()
// 		assert.Equal(t, value, expectedValue)
// 		assert.Equal(t, priority, expectedPriority)
// 	}
// }

// func TestMinPQueuePushAndPop_concurrently_protects_min_order(t *testing.T) {
// 	pqueue := NewPQueue(MINPQ)
// 	pqueueSize := 100

// 	var wg sync.WaitGroup

// 	// Populate the test priority queue with dummy elements
// 	// in asc ordered.
// 	for i := 0; i < pqueueSize; i++ {
// 		wg.Add(1)

// 		go func(i int) {
// 			defer wg.Done()

// 			value := NewTestJob(strconv.Itoa(i), i)

// 			pqueue.Push(value)
// 		}(i)
// 	}

// 	wg.Wait()

// 	for i := 0; i < pqueueSize; i++ {
// 		expectedValue := NewTestJob(strconv.Itoa(i), i)
// 		expectedPriority := int64(i)

// 		// Avoiding testing arithmetics headaches by using the pop function directly
// 		value, priority := pqueue.Pop()
// 		assert.Equal(t, value, expectedValue)
// 		assert.Equal(t, priority, expectedPriority)
// 	}
// }

// func TestMaxPQueueHead_returns_max_element(t *testing.T) {
// 	pqueue := NewPQueue(MAXPQ)

// 	pqueue.Push(NewTestJob("1", 1))
// 	pqueue.Push(NewTestJob("2", 2))

// 	value, priority := pqueue.Head()

// 	// First element of the binary heap is always left empty, so container
// 	// size is the number of elements actually stored + 1
// 	assert.Equal(t, len(pqueue.items), 3)
// 	assert.Equal(t, value, NewTestJob("2", 2))
// 	assert.Equal(t, priority, 2)
// }

// func TestMinPQueueHead_returns_min_element(t *testing.T) {
// 	pqueue := NewPQueue(MINPQ)

// 	pqueue.Push(NewTestJob("1", 1))
// 	pqueue.Push(NewTestJob("2", 2))

// 	value, priority := pqueue.Head()

// 	// First element of the binary heap is always left empty, so container
// 	// size is the number of elements actually stored + 1
// 	assert.Equal(t, len(pqueue.items), 3)
// 	assert.Equal(t, value, NewTestJob("1", 1))
// 	assert.Equal(t, priority, 1)
// }

// // -------
// // Benchmark for priority queue
// // -------
// func enqueueRandomJob(num int) {
// 	pqueue := NewPQueue(MINPQ)

// 	var job *Job
// 	for i := 0; i < num; i++ {
// 		job = NewJob(strconv.Itoa(i))
// 		job.NextRunAt = now().Add(time.Duration(rand.Intn(100000000)) * time.Minute)

// 		pqueue.Push(job)
// 	}
// }

// func benchmarkPQueue(i int, b *testing.B) {
// 	for n := 0; n < b.N; n++ {
// 		enqueueRandomJob(i) // enqueue i random events
// 	}
// }

// func BenchmarkPQueue10(b *testing.B) {
// 	benchmarkPQueue(10, b)
// }

// func BenchmarkPQueue20(b *testing.B) {
// 	benchmarkPQueue(20, b)
// }

// func BenchmarkPQueue50(b *testing.B) {
// 	benchmarkPQueue(50, b)
// }

// func BenchmarkPQueue100(b *testing.B) {
// 	benchmarkPQueue(100, b)
// }

// func BenchmarkPQueue500(b *testing.B) {
// 	benchmarkPQueue(500, b)
// }

// func BenchmarkPQueue1000(b *testing.B) {
// 	benchmarkPQueue(1000, b)
// }
// func BenchmarkPQueue2000(b *testing.B) {
// 	benchmarkPQueue(2000, b)
// }

// func BenchmarkPQueue5000(b *testing.B) {
// 	benchmarkPQueue(5000, b)
// }

// func BenchmarkPQueue50000(b *testing.B) {
// 	benchmarkPQueue(50000, b)
// }

// func BenchmarkPQueue500000(b *testing.B) {
// 	benchmarkPQueue(500000, b)
// }

// func BenchmarkPQueue1000000(b *testing.B) {
// 	benchmarkPQueue(1000000, b)
// }

// func BenchmarkPQueue100000000(b *testing.B) {
// 	benchmarkPQueue(100000000, b)
// }
