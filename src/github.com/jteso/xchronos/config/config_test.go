package config

import (
	"log"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/jteso/testify/assert"
)

const testConfig = `

version = "0.1"
job_store = ["1.1.1.1", "2.2.2.2"] 

job "outbound_email_marketing" {
    trigger = {
        cron = "* * * * 0/10"
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
		JobStore: []string{"1.1.1.1", "2.2.2.2"},
		Jobs: []JobConfig{
			JobConfig{
				Name: "outbound_email_marketing",
				Trigger: TriggerConfig{
					Cron:           "* * * * 0/10",
					Max_Executions: "-1",
				},
				Exec: "~/bin/outbound-email.sh",
			},
			JobConfig{
				Name: "Partner_Feeds_ETL",
				Trigger: TriggerConfig{
					Cron:           "* * * * 0/30",
					Max_Executions: "1",
				},
				Exec: "~/bin/partner-etl.sh",
			},
		},
	}
	config, err := ParseConfig(testConfig)

	assert.Nil(t, err, "Problems parsing the config file")
	assert.Equal(t, config.Version, expected.Version)
	assert.Equal(t, config.JobStore, expected.JobStore)
	assert.True(t, len(config.Jobs) == len(expected.Jobs))
	assert.Equal(t, config.Jobs, expected.Jobs)

	log.Println(spew.Sdump(config))
	// if !reflect.DeepEqual(config.Jobs[1], expected.Jobs[1]) {
	// 	t.Error("Config structure differed from expectation")
	// }

}
