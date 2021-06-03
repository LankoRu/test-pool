package pool

import (
	"context"
	"example/pool/limiter"
	"example/pool/queue"
	"example/pool/task"
	"sync"
)

// Pool is pool of workers
type Pool struct {
	queue   queue.Queue
	limiter limiter.Limiter
}

func New(q queue.Queue, l limiter.Limiter) *Pool {
	return &Pool{queue: q, limiter: l}
}

// Run runs pool processing in background, returns channel with results.
// Channel closes only if context finished.
func (p *Pool) Run(ctx context.Context) <-chan task.Result {
	resultsChan := make(chan task.Result)

	go func() {
		wg := sync.WaitGroup{}

		for {
			select {
			case <-ctx.Done():
				close(resultsChan)
				return
			case <-p.limiter.Next():
				nextTask, err := p.queue.GetNext(ctx)
				if err != nil {
					p.limiter.Ready()
					continue
				}

				wg.Add(1)
				go func() {
					defer wg.Done()

					result, err := p.executeTask(ctx, nextTask)
					if err == nil {
						resultsChan <- result
					}
				}()
			}
		}
	}()

	return resultsChan
}

func (p *Pool) executeTask(ctx context.Context, task task.Task) (task.Result, error) {
	defer p.limiter.Ready()

	return task.Execute(ctx)
}
