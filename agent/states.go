package agent

import(
	"github.com/jteso/xchronos/task"
)
// agentStateFn represents the state of an agent as a function
// that returns the next state
type handleStateFn func(*Agent) handleStateFn

func startStateFn(agent *Agent) handleStateFn {
	agent.changeState("STARTING_STATE")
	agent.connectEtcdCluster()

	return candidateStateFn
}

func candidateStateFn(agent *Agent) handleStateFn {
	agent.changeState("CANDIDATE_STATE")
	leader, err := agent.runForLeader()

	if err != nil {
		agent.lastError = err
		return errorStateFn
	}
	if leader {
		return leaderStateFn
	}

	return supporterStateFn

}

func leaderStateFn(agent *Agent) handleStateFn {
//	defer func(){
//		agent.StopPublishJobOffers()
//		agent.StopWatchForJobOffers()
//
//	}()
	agent.changeState("LEADER_STATE")

	// Advertise and renew the roles of this agent in etcd
	leaderTask := agent.advertiseAndRenewLeaderRoleT()
	executorTask := agent.advertiseAndRenewExecutorRoleT()

	// Tasks to be done as a scheduler leader
	//agent.SetupScheduler()
	jobPublisherTask := agent.publishJobOffersT()

	// Tasks to be done as an executor
	jobExecutorTask := agent.watchForJobOffersT()

	// Keep the lights on as an agent
	// agent.WatchForNewLeaderElection

	for {
		select {
		case err := <-task.FirstError(
							agent.ListenUserCancelTask(),
							leaderTask,
							jobPublisherTask,
							executorTask,
							jobExecutorTask):
			agent.lastError = err
			return errorStateFn
		}
	}

//		case <- agent.WatchForNewLeaderElection(): return candidateStateFn
//		case er := <- agent.ErrorReportedInTask(): return nil
//		case sig := <- agent.SignalReceived(): return nil
}

func supporterStateFn(agent *Agent) handleStateFn {
	agent.changeState("SUPPORTER_STATE")

	executorTask := agent.advertiseAndRenewExecutorRoleT()
	jobExecutorTask := agent.watchForJobOffersT()
	// [etcd leader watcher] return electForLeader
	// [etcd jobOffer watcher] competeForJobExecution
	for {
		select {
		case err := <-task.FirstError(
							agent.ListenUserCancelTask(),
							executorTask,
							jobExecutorTask):
			agent.lastError = err
			return errorStateFn
		}
	}
}


func errorStateFn(agent *Agent) handleStateFn {
	agent.logf("Entered in recovery mode due to error: %s", agent.lastError.Error())
	agent.stopTasks()
	return nil
}