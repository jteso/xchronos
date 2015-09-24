package config

import (
	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl"
)

type Config struct {
	Version  string
	JobStore []string
	Jobs     []JobConfig
}

type JobConfig struct {
	Name    string
	Trigger TriggerConfig
	Exec    string
	//OnError OnErrorConfig
}

type TriggerConfig struct {
	Cron           string
	Max_Executions string
}

// type OnErrorConfig struct {
// 	MaxRetries     string
// 	AbortOnFailure bool
// }

func ParseConfig(hclText string) (*Config, error) {
	config := &Config{}
	var errors *multierror.Error

	hclParseTree, err := hcl.Parse(hclText)
	//log.Println(spew.Sdump(hclParseTree))

	if err != nil {
		return nil, err
	}

	config.Version = parseVersion(hclParseTree, errors)
	config.JobStore = parseJobStore(hclParseTree, errors)
	config.Jobs = parseJobs(hclParseTree, errors)

	//log.Println(spew.Sdump(config.Jobs))

	return config, nil
}