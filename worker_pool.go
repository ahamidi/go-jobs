package jobs

// WorkerPool represents a collection of Workers
type WorkerPool struct {
	Workers []*Worker
}

// NewWorkerPool constructor
func NewWorkerPool(q Queue, f WorkerFunc, size int) (*WorkerPool, error) {
	return nil, nil
}

// Start runs the WorkerPool and will spin up the specified number of Workers
func (p *WorkerPool) Start() error {
	return nil
}

// Stop halts the WorkerPool and tears down all Workers
func (p *WorkerPool) Stop() error {
	return nil
}
