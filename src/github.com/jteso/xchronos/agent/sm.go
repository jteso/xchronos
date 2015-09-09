package agent

// implementation of the state machine that govern the behaviour of an agent
import (
	"github.com/jteso/task"
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
	agent.changeState("LEADER_STATE")

	leaderTask := agent.advertiseAndRenewLeaderRoleT()
	//TODO: agent.SetupScheduler()
	jobPublisherTask := agent.publishJobOffersT()

	executorTask := agent.advertiseAndRenewExecutorRoleT()
	jobExecutorTask := agent.watchForJobOffersT()
	watchNewLeaderTask := agent.watchForNewLeaderElectionT()

	for {
		select {
		case err := <-task.FirstError(
			leaderTask,
			jobPublisherTask,
			executorTask,
			jobExecutorTask):

			agent.lastError = err
			return errorStateFn
		case <-watchNewLeaderTask.ErrorChan():
			return candidateStateFn
		}
	}
}

func supporterStateFn(agent *Agent) handleStateFn {
	agent.changeState("SUPPORTER_STATE")

	executorTask := agent.advertiseAndRenewExecutorRoleT()
	jobExecutorTask := agent.watchForJobOffersT()
	watchNewLeaderTask := agent.watchForNewLeaderElectionT()

	for {
		select {
		case err := <-task.FirstError(
			executorTask,
			jobExecutorTask):

			agent.lastError = err
			return errorStateFn
		case <-watchNewLeaderTask.ErrorChan():
			return candidateStateFn

		}
	}
}

func errorStateFn(agent *Agent) handleStateFn {
	agent.changeState("RECOVERY_MODE_STATE")
	agent.logf("Error: %s", agent.lastError.Error())
	agent.stopTasks()
	return nil
}
