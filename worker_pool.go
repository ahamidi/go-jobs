package jobs

import (
	"sync"
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
		wg:      sync.WaitGroup{},
	}

	for i := 0; i < size; i++ {
		wp.Workers[i] = newWorker(q, fn, &wp.wg)
		wp.wg.Add(1)
	}
	return wp, nil
}

// Run the WorkerPool and will spin up the specified number of Workers
func (p *WorkerPool) Run() {
	for _, w := range p.Workers {
		go w.do()
	}
	p.wg.Wait()
}

// Stop halts the WorkerPool and tears down all Workers
func (p *WorkerPool) Stop() {
	p.stop <- struct{}{}
}
