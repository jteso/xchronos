package config

import (
	"log"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

const testConfig = `
version = "0.1"
job_store = "etcd" 

job "outbound_email_marketing" {
    trigger = {
        cron = "* * * * 0/30"
        max_executions = "-1"
    }
    exec = "~/bin/outbound-email.sh"
}

job "Partner_Feeds_ETL" {
    trigger = {
        cron = "* * * * 0/30"
        max_executions = "1"
    }
    exec = "~/bin/partner-etl.sh"
}

`

func TestConfigParsing(t *testing.T) {
	expected := &Config{
		Version:  "0.1",
		JobStore: "etcd",
		Jobs: []JobConfig{
			JobConfig{
				Name: "outbound_email_marketing",
				Trigger: TriggerConfig{
					Cron:          "* * * * 0/30",
					MaxExecutions: "-1",
				},
				Exec: "~/bin/outbound-email.sh",
			},
			JobConfig{
				Name: "Partner_Feeds_ETL",
				Trigger: TriggerConfig{
					Cron:          "* * * * 0/30",
					MaxExecutions: "-1",
				},
				Exec: "~/bin/partner-etl.sh",
			},
		},
	}
	config, err := ParseConfig(testConfig)
	if err != nil {
		log.Printf("*****\nError: %s\n", err.Error())
		t.Error(err)
	}

	if config.Version != expected.Version {
		t.Error("^^ Error here ^^")
	}
	if config.JobStore != expected.JobStore {
		t.Error("^^ Error here^^")
	}
	if len(config.Jobs) != len(expected.Jobs) {
		t.Error("Parsed incorrectly amount of jobs")
	}

	log.Println(spew.Sdump(config))
	// if !reflect.DeepEqual(config.Jobs[1], expected.Jobs[1]) {
	// 	t.Error("Config structure differed from expectation")
	// }

}
