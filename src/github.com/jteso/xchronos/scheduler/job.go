// A graduate student, Robert Brown, reviewing this article, recognized the parallel between cron
// and discrete event simulators, and created an implementation of the Franta-Maly event list manager for experimentation.

// The algorithm used by this cron is as follows:
// On start-up, look for a file named .crontab in the home directories of all account holders.
// For each crontab file found, determine the next time in the future that each command is to be run.
// Place those commands on the Franta-Maly event list with their corresponding time and their "five field" time specifier.
// Enter main loop:
//       1. Examine the task entry at the head of the queue, compute
//           how far in the future it is to be run.
//       2. Sleep for that period of time.
//       3. On awakening and after verifying the correct time, execute
//           the task at the head of the queue (in background) with the
//           privileges of the user who created it.
//       4. Determine the next time in the future to run this command
//           and place it back on the event list at that time value.

// The daemon would respond to SIGHUP signals to rescan modified crontab files and would schedule special "wake up events"
// on the hour and half hour to look for modified crontab files.

package scheduler

import (
	"bytes"
	"encoding/gob"
	"time"
)

// Any job will have to implement this interface
type Runnable interface {
	Run() error
	GetNextRunAt() time.Time
}

func EmptyJob() *Job {
	return &Job{}
}

type Job struct {
	Id        string
	NextRunAt time.Time
}

func NewJob(id string) *Job {
	return &Job{
		Id: id,
	}
}

func NewTestJob(id string, due int) *Job {
	return &Job{
		Id:        id,
		NextRunAt: time.Now().Add(time.Duration(due) * time.Second),
	}
}

func (j *Job) GetNextRunAt() time.Time {
	return j.NextRunAt
}

// WaitSecs returns number of secs to wait until job is due to run
func (j *Job) WaitSecs() float64 {
	secs := j.GetNextRunAt().Sub(time.Now()).Seconds()
	if secs < 0 {
		return 0
	}
	return secs
}

// Bytes returns the byte representation of the Job.
func (j Job) Bytes() ([]byte, error) {
	buff := new(bytes.Buffer)
	enc := gob.NewEncoder(buff)
	err := enc.Encode(j)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

// NewFromBytes returns a Job instance from a byte representation.
func NewFromBytes(b []byte) (*Job, error) {
	j := &Job{}

	buf := bytes.NewBuffer(b)
	err := gob.NewDecoder(buf).Decode(j)
	if err != nil {
		return nil, err
	}

	return j, nil
}
