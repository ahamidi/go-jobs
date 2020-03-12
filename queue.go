package jobs

// Queue is a logical collection of Jobs
type Queue struct {
	name string
	Database
}

// NewQueue constructor
func NewQueue(name string, db Database) *Queue {
	return &Queue{name, db}
}

// Next returns the next eligible Job in this Queue
func (q *Queue) Next() (*Job, error) {
	return q.Database.GetNextJob(q.name)
}

// Enqueue adds the given Job to the Queue
func (q *Queue) Enqueue(j *Job) (int, error) {
	return q.Database.EnqueueJob(q.name, j)
}
