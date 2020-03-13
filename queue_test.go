package jobs

import (
	"log"
	"testing"
)

func TestQueueNext(t *testing.T) {
	q, err := newQueue("queue")
	if err != nil {
		t.Errorf("expect no error, got %+v", err)
	}

	payload := map[string]string{
		"message": "hello",
	}
	j := NewJob(payload)
	id, err := q.Enqueue(j)
	if err != nil {
		t.Errorf("expect no error, got %+v", err)
	}

	job, err := q.Next()
	if err != nil {
		t.Errorf("expect no error, got %+v", err)
	}

	if job.ID != id {
		t.Errorf("expected id %d, go %d", id, job.ID)
	}

	t.Cleanup(func() {
		truncateJobsTable(q.DB)
	})
}

func BenchmarkQueueEnqueue(b *testing.B) {
	q, err := newQueue("benchmark")
	if err != nil {
		log.Fatal(err)
	}

	payload := map[string]string{
		"message": "hello",
	}
	for n := 0; n < b.N; n++ {
		j := NewJob(payload)

		_, err := q.Enqueue(j)
		if err != nil {
			log.Fatal(err)
		}
	}

	b.Cleanup(func() {
		truncateJobsTable(q.DB)
	})
}

func BenchmarkQueueNext(b *testing.B) {
	q, err := newQueue("benchmark")
	if err != nil {
		log.Fatal(err)
	}

	// setup
	payload := map[string]string{
		"message": "hello",
	}
	for i := 1; i < 10000; i++ {
		j := NewJob(payload)
		_, err := q.Enqueue(j)
		if err != nil {
			log.Fatal(err)
		}
	}
	b.ResetTimer()

	// benchmark
	for n := 0; n < b.N; n++ {
		job, err := q.Next()
		if err != nil {
			log.Fatal(err)
		}

		job.Complete(true, nil)
	}

	b.Cleanup(func() {
		truncateJobsTable(q.DB)
	})
}

func newQueue(name string) (*Queue, error) {
	return NewQueue(name, testDBURL())

}
