// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.

// This package contains some nice utilities to handle tasks, defining
// task as an arbitrary golang function that returns an error

// Example 1
// t := task.New("test", func() error {
//	fmt.Println("Hello World")
//	return nil
//})
// t.RunEvery(time.Second * 5)
//
// // To cancel it
// t.Stop()
//

package task

import (
	"errors"
	"time"
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

	onStopFn func()
}

func New(id string, fn func() error) *Task {
	return &Task{
		Id:       id,
		fn:       fn,
		errC:     make(chan error, 1),
		stopC:    make(chan chan struct{}, 1),
		onStopFn: func() {},
	}
}

// For testing purposes only
func NewDummy() *Task {
	return New("dummy", func() error {
		return nil
	})
}

// Stop gracefully a running task. i.e. this function will stop any timer or ticker may be running
// Keep in mind that this function will no close any channels may have been opened inside the task.
// See `OnStopFn()`
func (t *Task) Stop() {
	ackC := make(chan struct{}, 1)
	t.stopC <- ackC
	<-ackC
	t.onStopFn()
}

// Callback invoked right after the task has ack-ed the cancelation signal. It may be use to close any channels
// may have been opened inside the task.
func (t *Task) OnStopFn(callback func()) {
	t.onStopFn = callback
}

// implements the ErrorChanReader
func (t *Task) ErrorChan() chan error {
	return t.errC
}

func (t *Task) RunOnce() *Task {
	go func() {
		ackC := <-t.stopC
		ackC <- struct{}{}
		t.errC <- ErrUserCanceled
	}()

	go func() {
		t.errC <- t.fn()
	}()

	return t
}

func (t *Task) RunEvery(dur time.Duration) *Task {
	tkr := time.NewTicker(dur)
	go func() {

		for {
			select {
			case ackC := <-t.stopC:
				tkr.Stop()
				ackC <- struct{}{}
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

// === collection of tasks ===

// Given a number of tasks `ErrorChanReader`-ables, this function will emit the
// first error reported by any of them
func FirstError(tasks ...ErrorChanReader) chan error {
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
