package task

import (
	"time"
	"errors"
)

var (
	ErrUserCanceled = errors.New("Task canceled by user")
)

type ErrorChanReader interface {
	ErrorChan() chan error
}

type Task struct {
	Id string
	// The actual business function of a task
	fn func() error
	// Channel where all errors found by a task will be reported
	// if nil, the task finished with no errors
	errC chan error
	// Stop gracefully the running task. Do not use directly. Use the `Cancel` operation instead
	stopC chan chan struct{}

	// Force to stop a running task.
	killC chan struct{}

	postStopFn func()
}

func New(id string, fn func() error) *Task {
	return &Task{
		Id: id,
		fn: fn,
		errC: make(chan error, 1),
		stopC: make(chan chan struct{}, 1),
		postStopFn: func(){},
	}
}

// For testing purposes only
func NewDummy() *Task {
	return New("dummy", func() error {
		return nil
	})
}

// Stop gracefully a running task.
// Keep in mind that if task is scheduled to run once, this stop will have
// no effect
func (t *Task) Stop() {
	ackC := make(chan struct{}, 1)
	t.stopC <- ackC
	<- ackC
}

func (t *Task) PostStopFn(callback func()) {
	t.postStopFn = callback
}

// implements the ErrorChanReader
func (t *Task) ErrorChan() chan error{
	return t.errC
}


func (t *Task) RunOnce() *Task{
	go func() {
		ackC:= <- t.stopC
		ackC <- struct{}{}
		t.postStopFn()
		t.errC <- ErrUserCanceled
	}()

	go func(){
		t.errC <- t.fn()
	}()

	return t
}

func (t *Task) RunEvery(dur time.Duration) *Task{
	tkr := time.NewTicker(dur)
	go func() {

		for {
			select {
			case ackC := <- t.stopC:
				tkr.Stop()
				ackC <- struct{}{}
				t.postStopFn()
				t.errC <- ErrUserCanceled
				break
			case <-tkr.C:
				if err := t.fn(); err != nil {
					t.errC <- err
				}
			}
		}
	}()
	return t
}


// === collection of tasks

// Given a number of tasks `ErrorChanReader`-ables, this function will return the
// first error reported by any of them
func FirstError(tasks ...ErrorChanReader) chan error{
	firstErrC := make(chan error, 1)
	errReader := func(task ErrorChanReader, reportErrC chan error) {
		err := <-task.ErrorChan()
		if err != nil {
			// error found, lets report it
			reportErrC <- err
		}

	}

	for _, t := range tasks {
		go errReader(t, firstErrC)
	}

	return firstErrC
}
