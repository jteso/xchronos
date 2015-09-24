package errors

import "errors"

var (
	ErrNoSchedulerDetected = errors.New("cluster: No scheduler has been registered")
)
