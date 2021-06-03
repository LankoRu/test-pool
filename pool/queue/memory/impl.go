package memory

import (
	"context"
	"example/pool/queue"
	"example/pool/task"
	"sync"
)

type memoryQueue struct {
	sync.RWMutex
	tasks []task.Task
}

func New() queue.Queue {
	return &memoryQueue{tasks: make([]task.Task, 0, 100)}
}

func (q *memoryQueue) Add(_ context.Context, tasks ...task.Task) error {
	q.Lock()
	defer q.Unlock()

	q.tasks = append(q.tasks, tasks...)

	return nil
}

func (q *memoryQueue) GetNext(_ context.Context) (task.Task, error) {
	q.RLock()
	defer q.RUnlock()

	if len(q.tasks) == 0 {
		return nil, queue.ErrEmptyQueue
	}

	t := q.tasks[0]
	q.tasks = q.tasks[1:]

	return t, nil
}
