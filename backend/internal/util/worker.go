package util

import (
	"sync"

	"github.com/sqot0/packsmith/backend/internal/logger"
)

func WorkerPool[T any, R any](jobs <-chan T, fn func(T) R, numWorkers int) <-chan R {
	logger.Log.Printf("Starting worker pool with %d workers", numWorkers)
	out := make(chan R)
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				out <- fn(job)
			}
		}()
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
