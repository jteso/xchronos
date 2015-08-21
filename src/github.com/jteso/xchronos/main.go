package main

import (
	"flag"
	"os"
	"runtime"

	"github.com/jteso/xchronos/agent"

	"fmt"
	"time"
)

const (
	FLAG_ETCD_NODES = "etcd-nodes"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	os.Exit(realMain())
}

func parseFlags() map[string]string {
	result := make(map[string]string)

	etcdNodes := flag.String(FLAG_ETCD_NODES, "", "Comma separated list of etcd nodes")
	flag.Parse()

	result[FLAG_ETCD_NODES] = *etcdNodes
	return result
}

func realMain() int {
	flags := parseFlags()

	a1 := agent.New("agent_1", []string{flags[FLAG_ETCD_NODES]}, true)
	a2 := agent.New("agent_2", []string{flags[FLAG_ETCD_NODES]}, true)

	go a1.Run()
	go a2.Run()

	time.Sleep(1000 * time.Second)
	a1.Stop()
	a2.Stop()

	fmt.Println("System halted successfully :)")
	return 0
}
