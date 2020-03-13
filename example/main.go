package main

import (
	"errors"
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
		j := jobs.NewJob(payload)
		j.Retries = 10
		_, err := q.Enqueue(j)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(1000 * time.Millisecond)
	}

}

func counter(payload interface{}) error {
	p := payload.(map[string]interface{})
	log.Println("Processed value", p["value"].(float64))
	if p["value"].(float64) > 80 {
		return errors.New("value too high")
	}
	return nil
}
