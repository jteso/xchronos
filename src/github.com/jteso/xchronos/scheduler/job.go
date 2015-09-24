package scheduler

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"log"
	"os/exec"
	"strconv"
	"time"
	"fmt"
	"strings"

	"github.com/jteso/xchronos/config"
	"github.com/robfig/cron"
)

// Any job will have to implement this interface
type Runnable interface {
	Run() error
	GetNextRunAt() time.Time
	// If exhausted the number of max retries
	Failed() bool
}

func EmptyJob() *Job {
	return &Job{}
}

type Job struct {
	Id string

	// Command to run
	Exec string `json:"exec"`

	// Is this job disabled?
	Disabled bool `json:"disabled"`

	// Next time to run the job
	NextRunAt time.Time `json:"next_run_at`

	// Remaining executions, (-1) if infinite job
	TimesToRepeat    int64 `json:"times_to_repeat"`
	TimeToRepeatLeft int64 `json:"times_to_repeat_left"`

	Retries        uint `json:"retries"`
	CurrentRetries uint `json:"current_retries"`

	Trigger cron.Schedule
	// TODO(javier): add some stats tracking here
	// successCount, errorcounts, runattempt times,...
}

func NewJob(
	id string,
	exec string,
	disabled bool,
	timesToRepeat int64,
	retries uint,
	cronExp string) *Job {
	trigger, _ := cron.Parse(cronExp) // TODO(javier): error handling here

	return &Job{
		Id:               id,
		Exec:             exec,
		Disabled:         disabled,
		TimesToRepeat:    timesToRepeat,
		TimeToRepeatLeft: timesToRepeat,
		Retries:          retries - 1,
		CurrentRetries:   0,
		Trigger:          trigger,
		NextRunAt:        trigger.Next(time.Now().Add(-1 * time.Second)),
	}
}

func NewJobFromConfig(config config.JobConfig) *Job {
	timesToRepeat, _ := strconv.ParseInt(config.Trigger.Max_Executions, 10, 64)
	return NewJob(config.Name,
		config.Exec,
		false,
		timesToRepeat,
		3, // TODO(javier): FIXME
		config.Trigger.Cron)
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

// WaitSecs returns number of secs to wait until job is due for execution
func (j *Job) WaitSecs() float64 {
	secs := j.GetNextRunAt().Sub(time.Now()).Seconds()
	if secs < 0 {
		return 0
	}
	return secs
}

// Fire the job, this method does not check whether the job is due or not
// for execution. Hence, an external scheduler is required. See `scheduler.go`
func (j Job) Run() error {
	head, parts := splitCommand(j.Exec)
	out, err := exec.Command(head, parts...).Output()
	if err != nil {
			log.Print("Error found while executing job: %s due to %s", j.Id, err.Error())
		j.CurrentRetries ++
		//TODO(javier): add stats here
	} else{
		log.Printf("job.output ==> %s", out)
	}
	return err
}

func (j *Job) Failed() bool {
	return j.CurrentRetries > j.Retries
}

func (j Job) ToString() string {
	return fmt.Sprintf("[job: %s]", j.Id)
}

// Bytes returns the byte representation of the Job.
func (j Job) EncodeToString() (string, error) {
	buff := new(bytes.Buffer)
	gob.Register(&cron.SpecSchedule{})

	enc := gob.NewEncoder(buff)
	err := enc.Encode(j)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buff.Bytes()), nil
}

// NewFromBytes returns a Job instance from a byte representation.
func DecodeFromString(job string) (*Job, error) {
	jobAsBytes, _ := base64.StdEncoding.DecodeString(job)
	j := &Job{}
	gob.Register(&cron.SpecSchedule{})
	buf := bytes.NewBuffer(jobAsBytes)
	err := gob.NewDecoder(buf).Decode(j)
	if err != nil {
		return nil, err
	}

	return j, nil
}

func splitCommand(exec string) (head string, parts []string) {
	parts = strings.Fields(exec)
	head = parts[0]
	parts = parts[1:len(parts)]
	return head, parts
}
