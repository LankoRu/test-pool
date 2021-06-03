package queue

import (
	"context"
	"example/pool/task"
)

// Queue is interface for task storage
type Queue interface {
	Add(ctx context.Context, tasks ...task.Task) error
	GetNext(ctx context.Context) (task.Task, error)
}
