package main

import (
	"context"
	"example/pool"
	"example/pool/limiter/channel"
	"example/pool/queue/memory"
	"fmt"
	"log"
	"sync"
	"time"
)

func runInputGenerator() chan string {
	ch := make(chan string)

	go func() {
		for i := 0; i < 100; i++ {
			ch <- fmt.Sprintf("https://some-host/file/%d", i)
		}
		close(ch)
	}()

	return ch
}

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)

	wg := sync.WaitGroup{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	queue := memory.New()
	pool := pool.New(queue, channel.New(5))
	input := runInputGenerator()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for url := range input {
			err := queue.Add(ctx, NewDownloadTask(url))
			if err != nil {
				log.Printf("adding to queue failed - %v\n", err)
				continue
			}
		}
	}()

	resultsChan := pool.Run(ctx)
	for r := range resultsChan {
		log.Printf("result - %v\n", r)
	}

	wg.Wait()
}
