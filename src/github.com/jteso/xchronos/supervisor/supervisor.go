package supervisor

// TODO(javier): To provide different build-in restart strategies:
// - Exponential backoff supervision strategy: restart a failed job, each time with a growing time delay in between.
import(
	"errors"
	"github.com/jteso/xchronos/scheduler"
	_ "github.com/davecgh/go-spew/spew"
)
var(
	ErrExhaustedRetries = errors.New("Exceeded the number of max retries for this job")
)

type Supervisor struct {
}

func (s *Supervisor) RunIt(job *scheduler.Job) error {
	var err error
	//log.Println(spew.Sdump(job))
	for !job.Failed() {
		err = job.Run()
		if err == nil {
			return nil
		}
	}
	return ErrExhaustedRetries
}

