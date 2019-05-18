package jobs

import "time"

type Job struct {
	ID      int64       `json:"id"`
	Retries int         `json:"retries,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
	State   state       `json:"state"`
	Success bool        `json:"success,omitempty"`
	Error   error       `json:"error,omitempty"`

	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CompletedAt time.Time `json:"completed_at"`
}

func NewJob(payload interface{}) *Job {
	return &Job{
		ID:        0,
		Retries:   0,
		Payload:   payload,
		State:     New,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}
