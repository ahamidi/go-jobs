package jobs

type Database interface {
	FindJob(int64) (*Job, error)
	EnqueueJob(*Job) (int64, error)
	GetNextJob() (*Job, error)
	ListJobs(bool, int, int) ([]*Job, error)
	PendingJobs() int
	MarkJobAsCompleted(*Job) error
	MarkJobAsCompletedID(int64) error
	PauseJob(*Job) error
	PauseJobID(int64) error
}
