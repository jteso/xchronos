package scheduler

import (
	"math/rand"
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestMaxPQueue_init(t *testing.T) {
	pqueue := NewPQueue(MAXPQ)

	assert(
		t,
		len(pqueue.items) == 1,
		"len(pqueue.items) == %d; want %d", len(pqueue.items), 1,
	)

	assert(
		t,
		pqueue.Size() == 0,
		"pqueue.Size() = %d; want %d", pqueue.Size(), 0,
	)

	assert(
		t,
		pqueue.items[0] == nil,
		"pqueue.items[0] = %v; want %v", pqueue.items[0], nil,
	)

	assert(
		t,
		reflect.ValueOf(pqueue.comparator).Pointer() == reflect.ValueOf(max).Pointer(),
		"pqueue.comparator != max",
	)
}

func TestMinPQueue_init(t *testing.T) {
	pqueue := NewPQueue(MINPQ)

	assert(
		t,
		len(pqueue.items) == 1,
		"len(pqueue.items) = %d; want %d", len(pqueue.items), 1,
	)

	assert(
		t,
		pqueue.Size() == 0,
		"pqueue.Size() = %d; want %d", pqueue.Size(), 0,
	)

	assert(
		t,
		pqueue.items[0] == nil,
		"pqueue.items[0] = %v; want %v", pqueue.items[0], nil,
	)

	assert(
		t,
		reflect.ValueOf(pqueue.comparator).Pointer() == reflect.ValueOf(min).Pointer(),
		"pqueue.comparator != min",
	)
}

func TestMaxPQueuePushAndPop_protects_max_order(t *testing.T) {
	pqueue := NewPQueue(MAXPQ)
	pqueueSize := 100

	// Populate the test priority queue with dummy elements
	// in asc ordered.
	for i := 0; i < pqueueSize; i++ {
		value := NewJob(strconv.Itoa(i))
		priority := int64(i)

		pqueue.Push(value, priority)
	}

	containerIndex := 1 // binary heap are 1 indexed
	for i := 99; i >= 0; i-- {
		expectedValue := NewJob(strconv.Itoa(i))
		expectedPriority := int64(i)

		// Avoiding testing arithmetics headaches by using the pop function directly
		value, priority := pqueue.Pop()
		assert(
			t,
			reflect.DeepEqual(value, expectedValue),
			"value = %v; want %v", containerIndex, value, expectedValue,
		)
		assert(
			t,
			priority == expectedPriority,
			"priority = %v; want %v", containerIndex, priority, expectedValue,
		)

		containerIndex++
	}
}

func TestMaxPQueuePushAndPop_concurrently_protects_max_order(t *testing.T) {
	var wg sync.WaitGroup

	pqueue := NewPQueue(MAXPQ)
	pqueueSize := 100

	// Populate the test priority queue with dummy elements
	// in asc ordered.
	for i := 0; i < pqueueSize; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			value := NewJob(strconv.Itoa(i))
			priority := int64(i)
			pqueue.Push(value, priority)
		}(i)
	}

	wg.Wait()

	containerIndex := 1 // binary heap are 1 indexed
	for i := 99; i >= 0; i-- {
		expectedValue := NewJob(strconv.Itoa(i))
		expectedPriority := int64(i)

		// Avoiding testing arithmetics headaches by using the pop function directly
		value, priority := pqueue.Pop()
		assert(
			t,
			reflect.DeepEqual(value, expectedValue),
			"value = %v; want %v", containerIndex, value, expectedValue,
		)
		assert(
			t,
			priority == expectedPriority,
			"priority = %v; want %v", containerIndex, priority, expectedValue,
		)

		containerIndex++
	}
}

func TestMinPQueuePushAndPop_protects_min_order(t *testing.T) {
	pqueue := NewPQueue(MINPQ)
	pqueueSize := 100

	// Populate the test priority queue with dummy elements
	// in asc ordered.
	for i := 0; i < pqueueSize; i++ {
		value := NewJob(strconv.Itoa(i))
		priority := int64(i)

		pqueue.Push(value, priority)
	}

	for i := 0; i < pqueueSize; i++ {
		expectedValue := NewJob(strconv.Itoa(i))
		expectedPriority := int64(i)

		// Avoiding testing arithmetics headaches by using the pop function directly
		value, priority := pqueue.Pop()
		assert(
			t,
			reflect.DeepEqual(value, expectedValue),
			"value = %v; want %v", value, expectedValue,
		)
		assert(
			t,
			priority == expectedPriority,
			"priority = %v; want %v", priority, expectedValue,
		)
	}
}

func TestMinPQueuePushAndPop_concurrently_protects_min_order(t *testing.T) {
	pqueue := NewPQueue(MINPQ)
	pqueueSize := 100

	var wg sync.WaitGroup

	// Populate the test priority queue with dummy elements
	// in asc ordered.
	for i := 0; i < pqueueSize; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			value := NewJob(strconv.Itoa(i))
			priority := int64(i)

			pqueue.Push(value, priority)
		}(i)
	}

	wg.Wait()

	for i := 0; i < pqueueSize; i++ {
		expectedValue := NewJob(strconv.Itoa(i))
		expectedPriority := int64(i)

		// Avoiding testing arithmetics headaches by using the pop function directly
		value, priority := pqueue.Pop()
		assert(
			t,
			reflect.DeepEqual(value, expectedValue),
			"value = %v; want %v", value, expectedValue,
		)
		assert(
			t,
			priority == expectedPriority,
			"priority = %v; want %v", priority, expectedValue,
		)
	}
}

func TestMaxPQueueHead_returns_max_element(t *testing.T) {
	pqueue := NewPQueue(MAXPQ)

	pqueue.Push(NewJob("1"), 1)
	pqueue.Push(NewJob("2"), 2)

	value, priority := pqueue.Head()

	// First element of the binary heap is always left empty, so container
	// size is the number of elements actually stored + 1
	assert(t, len(pqueue.items) == 3, "len(pqueue.items) = %d; want %d", len(pqueue.items), 3)

	assertEqual(t, value, NewJob("2"), "pqueue.Head().value = %v; want %v", value, "2")
	assert(t, priority == 2, "pqueue.Head().priority = %d; want %d", priority, 2)
}

func TestMinPQueueHead_returns_min_element(t *testing.T) {
	pqueue := NewPQueue(MINPQ)

	pqueue.Push(NewJob("1"), 1)
	pqueue.Push(NewJob("2"), 2)

	value, priority := pqueue.Head()

	// First element of the binary heap is always left empty, so container
	// size is the number of elements actually stored + 1
	assert(t, len(pqueue.items) == 3, "len(pqueue.items) = %d; want %d", len(pqueue.items), 3)

	assertEqual(t, value, NewJob("1"), "pqueue.Head().value = %v; want %v", value, "1")
	assert(t, priority == 1, "pqueue.Head().priority = %d; want %d", priority, 1)
}

// -------
// Benchmark for priority queue
// -------
func enqueueRandomJob(num int) {
	pqueue := NewPQueue(MINPQ)

	var job *Job
	for i := 0; i < num; i++ {
		job = NewJob(strconv.Itoa(i))

		pqueue.Push(job, (now().Add(time.Duration(rand.Intn(100000000)) * time.Minute).UnixNano()))
	}
}

func benchmarkPQueue(i int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		enqueueRandomEvent(i) // enqueue i random events
	}
}

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

func BenchmarkPQueue500000(b *testing.B) {
	benchmarkPQueue(500000, b)
}

func BenchmarkPQueue1000000(b *testing.B) {
	benchmarkPQueue(1000000, b)
}

func BenchmarkPQueue100000000(b *testing.B) {
	benchmarkPQueue(100000000, b)
}
