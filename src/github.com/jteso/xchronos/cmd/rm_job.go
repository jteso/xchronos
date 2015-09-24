package cmd

import (
	"github.com/codegangsta/cli"
	_ "github.com/davecgh/go-spew/spew"
	"github.com/jteso/xchronos/cluster"
)

// Examples:
// xchronos rm --all
// xchronos rm --id=payment.batch
func RmJobCommand() cli.Command {
	return cli.Command{
		Name:  "rm",
		Usage: "schedule a batch of jobs for execution",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "all", Usage: "remove all jobs"},
		},
		Action: func(c *cli.Context) {
			rmJob(c)
		},
	}
}

// TODO(javier): Quick and dirty
func rmJob(c *cli.Context) {
	xchronosConfig := lookupXChronosConfig()
	proxy := cluster.NewEtcdClient(xchronosConfig.JobStore)
	proxy.UnregisterJobs()
}
