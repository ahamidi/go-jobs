package jobs

import (
	"encoding/json"
	"errors"
	"log"
	"time"
)

// MemDB is the actual in memory struct that holds the jobs
type MemDB struct {
	jobs []*Job
}

// NewMemDB returns new Memory Only job database
func NewMemDB() *MemDB {
	return &MemDB{
		[]*Job{},
	}
}

// FindJob returns a Job by ID
func (m *MemDB) FindJob(id int64) (*Job, error) {
	for _, j := range m.jobs {
		if j.ID == id {
			return j, nil
		}
	}
	return nil, errors.New("Not Found")
}

// EnqueueJob adds a Job for processing
func (m *MemDB) EnqueueJob(j *Job) (int64, error) {
	if j.ID > 0 {
		m.jobs[j.ID-1] = j
		return j.ID, nil
	}
	m.jobs = append(m.jobs, j)

	return int64(len(m.jobs)), nil
}

// GetNextJob returns the next job that is ready for processing. Specifically
// this returns a Job in the "New" state.
func (m *MemDB) GetNextJob() (*Job, error) {
	if len(m.jobs) == 0 {
		return nil, nil
	}
	var j *Job
	j, m.jobs = m.jobs[len(m.jobs)-1], m.jobs[:len(m.jobs)-1]
	return j, nil
}

// ListJobs returns multiple jobs, filtered by criteria specified
// - Unstarted: Return only "New" Jobs
// - Page: Specify the page to return
// - Limit: Limit the number of Jobs returned.
func (m *MemDB) ListJobs(unstarted bool, page int, limit int) ([]*Job, error) {
	return m.jobs, nil
}

func (m *MemDB) MarkJobAsCompleted(j *Job) error {
	return nil
}

func (m *MemDB) MarkJobAsCompletedID(id int64) error {
	j, err := m.FindJob(id)
	if err != nil {
		return err
	}
	j.State = Stopped
	j.UpdatedAt = time.Now().UTC()
	j.CompletedAt = time.Now().UTC()
	return nil
}

func (m *MemDB) PauseJob(j *Job) error {
	j.State = Paused
	j.UpdatedAt = time.Now().UTC()
	return nil
}

func (m *MemDB) PauseJobID(id int64) error {
	job, err := m.FindJob(id)
	if err != nil {
		return err
	}
	job.State = Paused
	job.UpdatedAt = time.Now().UTC()
	return nil
}

// Inspect outputs the in memory struct for convinient inspection
func (m *MemDB) Inspect() []byte {
	dbMap := map[string]interface{}{
		"jobs": m.jobs,
	}
	jsonDB, err := json.Marshal(dbMap)
	if err != nil {
		log.Println("Unable to inspect MemDB:", err)
		return nil
	}
	return jsonDB
}
