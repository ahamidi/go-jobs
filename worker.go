package jobs

import (
	"sync"
	"time"
)

// Worker struct contains reference to the Queue and the function to execture
// against jobs in that Queue
type Worker struct {
	Queue *Queue
	fn    WorkerFunc
	st    chan struct{}
	wg    *sync.WaitGroup
}

// WorkerFunc is the signature of Function needed to be passed into the Workers
type WorkerFunc func(payload interface{}) error

// NewWorker constructor
func newWorker(q *Queue, fn WorkerFunc, wg *sync.WaitGroup) *Worker {
	return &Worker{
		Queue: q,
		fn:    fn,
		st:    make(chan struct{}),
		wg:    wg,
	}
}

func (w *Worker) do() {
	for {
		select {
		case <-w.st:
			w.wg.Done()
			return
		default:
			j, err := w.Queue.Next()
			if err != nil {
				if err == NoPendingJobsError {
					time.Sleep(1 * time.Second)
				}
				continue
			}

			err = w.fn(j.Payload)
			if err != nil {
				j.Complete(false, err)
			}
			j.Complete(true, nil)
		}
	}
}

func (w *Worker) stop() {
	w.st <- struct{}{}
}
