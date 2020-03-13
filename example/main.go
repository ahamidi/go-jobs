package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/ahamidi/go-jobs"
)

func main() {
	// Get Database URL
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("$DATABASE_URL required")
	}

	// Create new queue
	q, err := jobs.NewQueue("counts", databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	go jobProducer(q)

	p, err := jobs.NewWorkerPool(q, counter, 1)
	if err != nil {
		log.Fatal(err)
	}

	p.Run()
}

func jobProducer(q *jobs.Queue) {

	for {
		payload := map[string]interface{}{
			"value": rand.Intn(100),
		}
		id, err := q.Enqueue(jobs.NewJob(payload))
		if err != nil {
			log.Fatal(err)
			log.Printf("Job Enqueued; id=%d", id)
		}
		time.Sleep(1 * time.Second)
	}

}

func counter(payload interface{}) error {
	p := payload.(map[string]interface{})
	log.Println("Processed value", p["value"].(float64))
	return nil
}
