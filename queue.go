package jobs

// Queue is a logical collection of Jobs
type Queue struct {
	name string
	DB   *Postgres
}

// NewQueue constructor
func NewQueue(name string, databaseURL string) (*Queue, error) {
	db, err := NewPG(databaseURL)
	if err != nil {
		return nil, err
	}
	return &Queue{name, db}, nil
}

// Next returns the next eligible Job in this Queue
func (q *Queue) Next() (*Job, error) {
	return q.DB.GetNextJob(q.name)
}

// Enqueue adds the given Job to the Queue
func (q *Queue) Enqueue(j *Job) (int, error) {
	return q.DB.EnqueueJob(q.name, j)
}
