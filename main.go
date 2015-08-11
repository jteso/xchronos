package main


import(
	"github.com/jteso/xchronos/agent"
	"runtime"
	"os"
	"time"
	"fmt"
)


func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	os.Exit(realMain())
}

// func realMain() int {
// 	client := etcd.NewClient([]string{"http://127.0.0.1:4001"})
// 	resp, err := client.Get("creds", false, false)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Printf("Current creds: %s: %s\n", resp.Node.Key, resp.Node.Value)
// 	watchChan := make(chan *etcd.Response)
// 	go client.Watch("/creds", 0 , false, watchChan, nil)
// 	log.Println("Waiting for an update...")
// 	r := <- watchChan
// 	log.Printf("Got updated creds: %s: %s \n", r.Node.Key, r.Node.Value)
// 	return 0
// }

func realMain() int{

	a1_stopCh := make(agent.AckableOpCh)
	a2_stopCh := make(agent.AckableOpCh)

	a1 := agent.New("agent_1", a1_stopCh)
	a2 := agent.New("agent_2", a2_stopCh)

	
	go a1.Run()
	go a2.Run()
	
	time.Sleep(5 * time.Second)
	
	// fmt.Println("Sending stop signal to agents...")

	// terminate_1 := agent.NewTerminateOp()
	// terminate_2 := agent.NewTerminateOp()

	// a1_stopCh <- terminate_1
	// a2_stopCh <- terminate_2

	
	// fmt.Println("Waiting for agents to stop gracefully")	
	

	// terminate_1.WaitForAck()
	// terminate_2.WaitForAck()


	fmt.Println("System halted successfully :)")
	return 0
}