package agent

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
		return agent.panic("Error while electing for leadership due to: %s\n", err.Error())
	}
	if leader {
		return leaderStateFn
	}

	return supporterStateFn

}

func leaderStateFn(agent *Agent) handleStateFn {
	agent.changeState("LEADER_STATE")
	execErrCh := agent.registerAsExecutor()
	pubJobErrCh := agent.publishJobOffers()

	// agent.PersistCronConf()
	// agent.PublishJobs(){
	// 			for job := agent.scheduler(cronConf) {
	//      		etcd.OfferJob(job, ROUND_ROBIN_POLICY)
	// 			}
	//
	// agent.AcceptJobs()
	// ready <- agent.ListenForChangeOfLeadership()
	// agent.StopPublishJobs()
	// agent.StopAcceptingJobs()
	for {
		select {
		case <- execErrCh: return nil
		case <- pubJobErrCh: return nil
		}
	}

}

func supporterStateFn(agent *Agent) handleStateFn {
	agent.changeState("SUPPORTER_STATE")
	watchJobsErrCh := agent.watchForJobOffers()
	// [etcd leader watcher] return electForLeader
	// [etcd jobOffer watcher] competeForJobExecution
	<- watchJobsErrCh
	return nil
}
