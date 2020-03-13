package jobs

import (
	"sync"
	"time"
)

// WorkerPool represents a collection of Workers
type WorkerPool struct {
	Workers []*Worker
	wg      sync.WaitGroup
	stop    chan struct{}
}

// NewWorkerPool constructor
func NewWorkerPool(q *Queue, fn WorkerFunc, size int) (*WorkerPool, error) {
	wp := &WorkerPool{
		Workers: make([]*Worker, size),
	}

	for i := 0; i < size; i++ {
		wp.Workers[i] = newWorker(q, fn)
		wp.wg.Add(1)
	}
	return wp, nil
}

// Run the WorkerPool and will spin up the specified number of Workers
func (p *WorkerPool) Run() {
	for {
		select {
		case <-p.stop:
			break
		default:
			for _, w := range p.Workers {
				go func(w *Worker) {
					err := w.do()
					if err == NoPendingJobsError {
						time.Sleep(1 * time.Second)
					}
				}(w)
			}
		}
	}
}

// Stop halts the WorkerPool and tears down all Workers
func (p *WorkerPool) Stop() {
	p.stop <- struct{}{}
}
