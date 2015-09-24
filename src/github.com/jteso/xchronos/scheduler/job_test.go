package scheduler

import (
	"log"
	"testing"
	"time"

	"github.com/jteso/testify/assert"
)

var cronsExp = map[string]string{
	"every_second": "* * * * * *",
	"every_minute": "0 * * * * *",
}

func TestJob(t *testing.T) {

	job := NewJob("test", "script.sh", false, -1, 3, cronsExp["every_minute"])
	log.Printf("> Now time : %s", time.Now())
	log.Printf("> Next time: %s", job.GetNextRunAt())
	log.Printf("> Wait time: %0.1f (secs)", job.WaitSecs())
}

func TestEncDecJob(t *testing.T) {
	job := NewJob("test", "script.sh", false, -1, 3, cronsExp["every_minute"])
	bytes, err := job.Bytes()

	assert.NoError(t, err)

	job2, err2 := NewFromBytes(bytes)
	assert.NoError(t, err2)

	assert.Equal(t, job, job2, "job before and after serialization is not the same")
}
