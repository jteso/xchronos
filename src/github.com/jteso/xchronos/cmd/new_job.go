package cmd

import (
	"fmt"

	"github.com/codegangsta/cli"
	_ "github.com/davecgh/go-spew/spew"
	"github.com/jteso/xchronos/cluster"
	"github.com/jteso/xchronos/scheduler"
)

// Examples:
// xchronos add .
// xchronos add payment.batch
func NewJobCommand() cli.Command {
	return cli.Command{
		Name:  "add",
		Usage: "schedule a batch of jobs for execution",
		Action: func(c *cli.Context) {
			addJob(c)
		},
	}
}

// TODO(javier): Quick and dirty
func addJob(c *cli.Context) {
	//fmt.Println(spew.Sdump(c))
	var location string
	if len(c.Args()) == 0 {
		// lets assume you want to load batch files in the current dir
		location = "."
	} else {
		location = c.Args()[0] //TODO(javier): fix this
	}
	xchronoConfig := lookupXChronosConfig()
	jobConfig := parseJobs(location)
	//log.Println(spew.Sdump(jobConfig[0].Jobs[0]))
	fmt.Printf("Found %d batch files \n", len(jobConfig))

	proxy := cluster.NewEtcdClient(xchronoConfig.JobStore)
	fmt.Printf("==> Registering job: %s ...", jobConfig[0].Jobs[0].Name)
	proxy.RegisterJob(scheduler.NewJobFromConfig(jobConfig[0].Jobs[0]))
}
