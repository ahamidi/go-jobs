package jobs_test

import (
	"testing"

	jobs "github.com/ahamidi/go-jobs"
)

func TestNewJob(t *testing.T) {
	payload := "test payload"
	j := jobs.NewJob(payload)

	if j.State != jobs.New {
		t.Error("Job should be in state \"new\"")
	}
}
