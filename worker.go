package jobs

// Worker struct contains reference to the Queue and the function to execture
// against jobs in that Queue
type Worker struct {
	Queue Queue
	fn    WorkerFunc
}

// WorkerFunc is the signature of Function needed to be passed into the Workers
type WorkerFunc func(payload interface{}) error

// NewWorker constructor
func NewWorker(q Queue) *Worker {
	return &Worker{
		Queue: q,
	}
}

func (w *Worker) do() error {
	j, err := w.Queue.Next()
	if err != nil {
		return err
	}

	return w.fn(j.Payload)
}
