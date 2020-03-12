package jobs

type WorkerPool struct {
	Workers []*Worker
}

func NewWorkerPool(q Queue, f WorkerFunc, size int) (*WorkerPool, error) {
	return nil, nil
}

func (p *WorkerPool) Start() error {
	return nil
}

func (p *WorkerPool) Stop() error {
	return nil
}
