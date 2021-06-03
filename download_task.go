package main

import (
	"context"
	"example/pool/task"
	"math/rand"
	"path"
	"time"
)

// DownloadTask is test task, emulates downloading of some data by url
type DownloadTask struct {
	url string
}

func NewDownloadTask(url string) *DownloadTask {
	return &DownloadTask{url: url}
}

func (d DownloadTask) Execute(ctx context.Context) (task.Result, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err() // TODO Custom action
	case <-time.After(time.Duration(rand.Int()%1000+100) * time.Millisecond):
		return &DownloadTaskResult{DownloadedSize: path.Base(d.url)}, nil
	}
}

// DownloadTaskResult is test task result
type DownloadTaskResult struct {
	DownloadedSize string
}

func (r *DownloadTaskResult) Value() interface{} {
	return r
}
