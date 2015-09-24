package agent

// implementation of the state machine that govern the behaviour of an agent
import "github.com/jteso/xchronos/errors"

// agentStateFn represents the state of an agent as a function
// that returns the next state
type handleStateFn func(*Agent) handleStateFn

func startStateFn(agent *Agent) handleStateFn {
	agent.NotifyState("_connecting")
	agent.ConnectCluster()

	return negotiateStateFn
}

func negotiateStateFn(agent *Agent) handleStateFn {
	agent.NotifyState("_pending")
	leader, err := agent.AttemptJobScheduler()

	if err != nil {
		agent.lastError = err
		return errorStateFn
	}
	if leader {
		return schedulerStateFn
	}

	return executorStateFn

}

func schedulerStateFn(agent *Agent) handleStateFn {
	agent.NotifyState("_scheduler")

	err := <-agent.ActAsJobScheduler()

	if err != nil && err == errors.ErrNoSchedulerDetected {
		return negotiateStateFn
	} else {
		return errorStateFn
	}

}

func executorStateFn(agent *Agent) handleStateFn {
	agent.NotifyState("_executor")

	err := <-agent.ActAsJobExecutor()

	if err != nil && err == errors.ErrNoSchedulerDetected {
		return negotiateStateFn
	} else {
		return errorStateFn
	}

}

func errorStateFn(agent *Agent) handleStateFn {
	agent.NotifyState("_recovery_mode")
	agent.Halt()
	return nil
}
