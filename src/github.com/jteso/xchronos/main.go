package main

import (
	"os"
	"runtime"

	"github.com/codegangsta/cli"
	"github.com/jteso/xchronos/cmd"
)

const (
	ETCD_NODES = "etcd-nodes"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	os.Exit(realMain())
}

func realMain() int {
	app := cli.NewApp()
	app.Name = "xchronos"
	app.Version = "0.1-Alpha"
	app.Usage = "Enabling your jobs to run in the cloud (public/private/hybrid)"
	app.Flags = []cli.Flag{
	//	cli.BoolFlag { Name: "debug", Usage: "output all cluster and agent activity"}
	}
	app.Commands = []cli.Command{
		cmd.RunAgentCommand(),
		//cmd.RunConsoleCommand(),
		cmd.NewJobCommand(),
		cmd.RmJobCommand(),
	}
	err := app.Run(os.Args)
	if err != nil {
		return 1
	}
	return 0
}
