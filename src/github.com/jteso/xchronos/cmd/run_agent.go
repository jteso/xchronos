package cmd

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/jteso/crypto"
	"github.com/jteso/xchronos/agent"
)

// <executable> agent (--no-scheduler | --no-executor)
func RunAgentCommand() cli.Command {
	return cli.Command{
		Name:  "agent",
		Usage: "run a xchronos agent in the local host",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "verbose", Usage: "increase all cluster and agent verbosity"},
		},
		Action: func(c *cli.Context) {
			runAgent(c)
		},
	}
}

func runAgent(c *cli.Context) {
	uuid, _ := crypto.GenerateUUID()
	config := lookupXChronosConfig()
	//fmt.Printf("job store: %+v\n", config.JobStore)
	//fmt.Printf("verbose: %t\n", c.Bool("verbose"))
	a := agent.New(uuid, config.JobStore, c.Bool("verbose"))
	fmt.Printf("==> Agent starting...\n")
	a.Run()
}
