package task

import (
	"log"
	"sync"
)

type TaskManager struct {
	tasks []*Task
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks: []*Task{},
	}
}

func (tm *TaskManager) RegisterTask(t *Task) {
	tm.tasks = append(tm.tasks, t)
}

func (t *TaskManager) Size() int {
	return len(t.tasks)
}

func (tm *TaskManager) StopAllTasks(doneC chan bool) {
	var wg sync.WaitGroup
	wg.Add(len(tm.tasks))

	for _, t := range tm.tasks {
		go func(t *Task) {
			log.Printf("Task: %s stopping...", t.Id)
			t.Stop()
			log.Printf("Task: %s stopped", t.Id)
			wg.Done()
		}(t)
	}

	log.Printf("Waiting for %d tasks to stop...", len(tm.tasks))
	wg.Wait()

	tm.tasks = tm.tasks[:0]

	doneC <- true
}
