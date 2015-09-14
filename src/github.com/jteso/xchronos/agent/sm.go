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
	agent.connectCluster()

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

	// tasks as scheduler
	schedulerRoleTask := agent.advertiseSchedulerRoleTask()
	runSchedulerTask := agent.runSchedulerTask()

	// tasks as job executor
	executorRoleTask := agent.advertiseExecutorRoleTask()
	offersListenerTask := agent.watchForJobOffersTask()
	schedulerFailureWatcherTask := agent.watchForSchedulerFailureTask()

	for {
		select {
		case err := <-task.FirstError(
			schedulerRoleTask,
			runSchedulerTask,
			executorRoleTask,
			offersListenerTask):

			agent.lastError = err
			return errorStateFn
		case <-schedulerFailureWatcherTask.ErrorChan():
			return candidateStateFn
		}
	}
}

func supporterStateFn(agent *Agent) handleStateFn {
	agent.changeState("SUPPORTER_STATE")

	executorRoleTask := agent.advertiseExecutorRoleTask()
	offersListenerTask := agent.watchForJobOffersTask()
	schedulerFailureWatcherTask := agent.watchForSchedulerFailureTask()

	for {
		select {
		case err := <-task.FirstError(
			executorRoleTask,
			offersListenerTask):

			agent.lastError = err
			return errorStateFn
		case <-schedulerFailureWatcherTask.ErrorChan():
			return candidateStateFn

		}
	}
}

func errorStateFn(agent *Agent) handleStateFn {
	agent.changeState("RECOVERY_MODE_STATE")
	agent.logf("Error: %s", agent.lastError.Error())
	agent.Stop()
	return nil
}
