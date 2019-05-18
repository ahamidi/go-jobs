package jobs

type queue struct {
	name string
}

func NewQueue(name string) *queue {
	return &queue{name}
}

func (q *queue) Next() *Job {
	return &Job{}
}
