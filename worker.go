package jobs

type Worker struct {
	Queue Queue
}

type WorkerFunc func(payload interface{}) error

func NewWorker(q Queue) *Worker {
	return &Worker{
		Queue: q,
	}
}

func (w *Worker) Do(fn WorkerFunc) error {
	j, err := w.Queue.Next()
	if err != nil {
		return nil
	}

	return fn(j.Payload)
}
