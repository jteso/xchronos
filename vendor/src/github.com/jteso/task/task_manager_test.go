// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.

package task

import (
	"testing"
	"time"
)

func TestTaskManager(t *testing.T) {
	runCounts := 0

	tm := NewTaskManager()

	// Task 1
	t1 := New("t1", func() error {
		runCounts++
		return nil
	})

	t1.RunEvery(time.Second * 1)
	tm.RegisterTask(t1)

	// Task 2
	t2 := New("t2", func() error {
		runCounts++
		return nil
	})

	t2.RunEvery(time.Second * 1)
	tm.RegisterTask(t2)

	// Stop all tasks
	doneC := make(chan bool, 1)
	tm.StopAllTasks(doneC)
	<-doneC

	// Validate that the task has been canceled by checking whether or not
	// the errorC has been populated
	var err error
	timeout := time.NewTimer(time.Second * 2)
	select {
	case err = <-t1.ErrorChan():
		if err != ErrUserCanceled {
			t.Errorf("Expected to receive a user cancelation error")
		}
		if tm.Size() != 0 {
			t.Errorf("Expected: <0> , Actual:<%d>", tm.Size())
		}
	case <-timeout.C:
		t.Errorf("Expected to receive a user cancelation error")
	}
}
