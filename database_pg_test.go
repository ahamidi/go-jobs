package jobs

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"
)

func TestEnqueueJob(t *testing.T) {
	db, err := pg()
	if err != nil {
		log.Fatal(err)
	}

	j := NewJob(map[string]interface{}{
		"message": "hello",
	})

	id, err := db.EnqueueJob("messages", j)
	if err != nil {
		t.Errorf("expect no error, got %+v", err)
	}

	if id == 0 {
		t.Errorf("expect id not 0, got %d", id)
	}

	t.Cleanup(func() {
		truncateJobsTable(db)
	})
}

func TestFindJob(t *testing.T) {
	db, err := pg()
	if err != nil {
		log.Fatal(err)
	}

	payload := map[string]string{
		"message": "hello",
	}
	j := NewJob(payload)

	id, _ := db.EnqueueJob("messages", j)

	job, err := db.FindJob(id)
	if err != nil {
		t.Errorf("expect no error, got %+v", err)
	}

	if m := job.Payload.(map[string]interface{})["message"].(string); m != payload["message"] {
		t.Errorf("expect payload message %s, got %s", payload["message"], m)
	}

	t.Cleanup(func() {
		truncateJobsTable(db)
	})
}

func TestGetNextJob(t *testing.T) {
	db, err := pg()
	if err != nil {
		log.Fatal(err)
	}

	// Add 2 jobs
	payload := map[string]string{
		"message": "hello",
	}
	j1 := NewJob(payload)
	j1.Retries = 3
	id1, _ := db.EnqueueJob("test", j1)
	j2 := NewJob(payload)
	j2.Retries = 3
	id2, _ := db.EnqueueJob("test", j2)

	// Should get first Job
	job1, err := db.GetNextJob("test")
	if err != nil {
		t.Errorf("expect no error, got %+v", err)
	}

	if job1.ID != id1 {
		t.Errorf("expect id %d, got %d", id1, job1.ID)
	}

	// Without "completing" first one, we should now get the second
	job2, err := db.GetNextJob("test")
	if err != nil {
		t.Errorf("expect no error, got %+v", err)
	}

	if job2.ID != id2 {
		t.Errorf("expect id %d, got %d", id2, job2.ID)
	}

	// Complete the first job
	err = job1.Complete(false, errors.New("timeout"))
	if err != nil {
		t.Errorf("expect no error, got %+v", err)
	}
	job3, err := db.GetNextJob("test")
	if err != nil {
		t.Errorf("expect no error, got %+v", err)
	}

	if job3.ID != id1 {
		t.Errorf("expect id %d, got %d", id1, job3.ID)
	}

	job2.Complete(true, nil)
	job3.Complete(true, nil)

	t.Cleanup(func() {
		truncateJobsTable(db)
	})
}

func pg() (*Postgres, error) {
	u := testDBURL()
	return NewPG(u)
}

func truncateJobsTable(db *Postgres) {
	db.Exec(context.Background(), "TRUNCATE TABLE jobs")
	//db.Exec(context.Background(), "delete from jobs")
}

func testDBURL() string {
	u := os.Getenv("TEST_DATABASE_URL")
	if u == "" {
		log.Fatal("no test database provided")
	}

	u = u + "&pool_max_conns=10"
	return u
}
