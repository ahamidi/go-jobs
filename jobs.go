package jobs

import (
	"time"
)

// Job represents a background Job
type Job struct {
	ID      int         `json:"id"`
	Retries int         `json:"retries,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
	State   state       `json:"state"`
	Success *bool       `json:"success,omitempty"`
	Error   *error      `json:"error,omitempty"`

	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at"`
	tx          Transaction
}

// NewJob constructor
func NewJob(payload interface{}) *Job {
	return &Job{
		Retries:   0,
		Payload:   payload,
		State:     New,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

// Complete marks the job as completed (either successfully or unsuccessfully
func (j *Job) Complete(success bool, e error) error {
	return j.tx.MarkJobAsCompleted(j, success, e)
}
