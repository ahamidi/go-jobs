package jobs

// Worker struct contains reference to the Queue and the function to execture
// against jobs in that Queue
type Worker struct {
	Queue *Queue
	fn    WorkerFunc
}

// WorkerFunc is the signature of Function needed to be passed into the Workers
type WorkerFunc func(payload interface{}) error

// NewWorker constructor
func newWorker(q *Queue, fn WorkerFunc) *Worker {
	return &Worker{
		Queue: q,
		fn:    fn,
	}
}

func (w *Worker) do() error {
	j, err := w.Queue.Next()
	if err != nil {
		return err
	}

	err = w.fn(j.Payload)
	if err != nil {
		j.Complete(false, err)
		return err
	}
	j.Complete(true, nil)
	return nil
}
