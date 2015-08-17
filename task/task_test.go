package task

import (
	"github.com/jteso/xchronos/task"
	"testing"
	"time"
)

func TestEvery(t *testing.T){
	runCounts := 0

	tsk := task.New(func() error {
		runCounts ++
		return nil
	})

	tsk.RunEvery(time.Second * 1)

	timeout := time.NewTimer(time.Second * 2)
	<- timeout.C

	t.Logf("RunCounts: %d\n", runCounts)

	// Assert that at this point, task should have run at least 4,5 times
	if runCounts != 2 {
		t.Errorf("Expected [%d] runs. Observed [%d] runs\n", 2, runCounts)
	}

}


func TestEveryCancel(t *testing.T){
	runCounts := 0

	tsk := task.New(func() error {
		runCounts ++
		return nil
	})

	tsk.RunEvery(time.Second * 1)

	timeKiller := time.NewTimer(time.Second * 2)
	<- timeKiller.C

	tsk.Stop()

	// Validate that the task has been canceled by checking whether or not
	// the errorC has been populated
	var err error
	timeout := time.NewTimer(time.Second * 1)
	select {
	case err = <- tsk.ErrC:
		if runCounts != 2{
			t.Errorf("Expected [%d] runs. Observed [%d] runs\n", 2, runCounts)
		}
	case <-timeout.C:
		t.Errorf("Expected to receive a user cancelation error")
	}

	t.Logf("Received error: %s", err.Error())
}

func TestFindFirstError(t *testing.T){

	tsk1 := task.NewDummy().RunEvery(time.Second * 1)
	tsk2 := task.NewDummy().RunEvery(time.Second * 1)

	// Stop one task after timeout
	go func() {
		timeout := time.NewTimer(time.Second * 2)
		<-timeout.C
		tsk1.Stop()
	}()

	firstErr:= <- task.FirstError(tsk1, tsk2)
	if firstErr.Error() != task.ErrUserCanceled.Error() {
		t.Errorf("Expected to get an user cancelation error")
	}

}

func TestKill(t *testing.T){
	count := 0
	tsk1 := task.New(func() error{
		for{
			t.Logf("Counting forever until someone kill me. Count: %d", count)
			count++
		}
		return nil
	})

	tsk1.RunOnce() // it will run forever

	// KILL the task after timeout
	go func() {
		timeout := time.NewTimer(time.Second * 2)
		<- timeout.C
		tsk1.Kill()
	}()

	err <- tsk1.ErrorChan()
	t.Logf("Task exited due to err: %s", err.Error())
}
