package jobs

// Database interface
type Database interface {
	FindJob(int) (*Job, error)
	EnqueueJob(string, *Job) (int, error)
	GetNextJob(string) (*Job, error)
	ListJobs(string, bool, int, int) ([]*Job, error)
	PendingJobs(string) (int, error)
	MarkJobAsCompleted(*Job, bool, error) error
}
