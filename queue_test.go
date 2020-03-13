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
}

func BenchmarkQueueEnqueue(b *testing.B) {
	q, err := newQueue("benchmark")
	if err != nil {
		log.Fatal(err)
	}

	truncateJobsTable(q.DB)

	for n := 0; n < b.N; n++ {
		payload := map[string]string{
			"message": "hello",
		}
		j := NewJob(payload)

		_, err := q.Enqueue(j)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func newQueue(name string) (*Queue, error) {
	return NewQueue(name, testDBURL())

}
