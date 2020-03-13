package jobs

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgxpool"
)

// Postgres handle
type Postgres struct {
	*pgxpool.Pool
}

// Transaction wrapper to support extending
type Transaction struct {
	pgx.Tx
}

// NewPG returns a verified Postgres struct
func NewPG(u string) (*Postgres, error) {
	pool, err := pgxpool.Connect(context.Background(), u)
	if err != nil {
		return nil, err
	}

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	if err = conn.Conn().Ping(context.Background()); err != nil {
		return nil, err
	}

	return &Postgres{pool}, nil
}

// FindJob returns a single Job with the given ID
func (db *Postgres) FindJob(id int) (*Job, error) {
	var retries int
	var payload interface{}
	var state state
	var success *bool
	var eString *string
	var createdAt, updatedAt time.Time
	var completedAt *time.Time

	err := db.QueryRow(context.Background(), "SELECT id, retries, payload, state, success, error, created_at, updated_at, completed_at FROM jobs WHERE id = $1", id).
		Scan(&id, &retries, &payload, &state, &success, &eString, &createdAt, &updatedAt, &completedAt)
	if err != nil {
		return nil, err
	}
	var e error
	if e != nil {
		e = errors.New(*eString)
	}

	return &Job{
		ID:          id,
		Retries:     retries,
		Payload:     payload,
		State:       state,
		Success:     success,
		Error:       &e,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		CompletedAt: completedAt,
	}, nil
}

// EnqueueJob enqueue the provided Job in the given Queue
func (db *Postgres) EnqueueJob(queue string, j *Job) (int, error) {
	query := `INSERT INTO jobs (queue, retries, payload) values ($1,$2,$3) RETURNING id;`

	var id int
	return id, db.QueryRow(context.Background(), query, queue, j.Retries, j.Payload).Scan(&id)
}

// GetNextJob returns the next unclaimed Job from the given Queue. The Job will
// be "claimed" and will not be visible by other callers unless it is "returned"
func (db *Postgres) GetNextJob(queue string) (*Job, error) {
	var id, retries int
	var payload interface{}
	var state state
	var success *bool
	var eString *string
	var createdAt, updatedAt time.Time
	var completedAt *time.Time

	tx, err := db.Begin(context.Background())
	if err != nil {
		return nil, err
	}

	err = tx.QueryRow(context.Background(), "SELECT id, retries, payload, state, success, error, created_at, updated_at, completed_at FROM jobs WHERE queue = $1 AND state = 0 ORDER BY updated_at ASC LIMIT 1 FOR UPDATE SKIP LOCKED", queue).
		Scan(&id, &retries, &payload, &state, &success, &eString, &createdAt, &updatedAt, &completedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, NoPendingJobsError
		}
		return nil, err
	}
	var e error
	if e != nil {
		e = errors.New(*eString)
	}

	return &Job{
		ID:          id,
		Retries:     retries,
		Payload:     payload,
		State:       Running,
		Success:     success,
		Error:       &e,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		CompletedAt: completedAt,
		tx:          Transaction{tx},
	}, nil
}

// ListJobs returns an array of Jobs from the given Queue, filtered by criteria specified
// - includeClaimed: Include Jobs that are currently being processed
// - page: Page to return
// - limit:  Limit the numbers of jobs returned
func (db *Postgres) ListJobs(queue string, includeClaimed bool, page int, limit int) ([]*Job, error) {
	return nil, nil
}

// PendingJobs returns a count of outstanding unclaimed Jobs in the given Queue
func (db *Postgres) PendingJobs(queue string) (int, error) {
	var c int
	err := db.QueryRow(context.Background(), "SELECT COUNT(id) WHERE queue = $1 AND state = 0;", queue).Scan(&c)
	if err != nil {
		return 0, err
	}
	return c, nil
}

// MarkJobAsCompleted Marks the given Job as having been completed and whether
// it was successful or failed
func (tx Transaction) MarkJobAsCompleted(j *Job, successful bool, e error) error {
	var s state
	if j.Retries == 0 {
		if successful {
			s = Completed
		} else {
			s = Failed
		}
	} else {
		j.Retries--
	}

	if e == nil {
		e = errors.New("no error")
	}

	ct, err := tx.Exec(context.Background(), "UPDATE jobs SET retries = $1, state = $2, success = $3, error = $4, updated_at = CURRENT_TIMESTAMP, completed_at = CURRENT_TIMESTAMP WHERE id = $5", j.Retries, s, successful, e.Error(), j.ID)
	if err != nil {
		return err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return errors.New("no jobs updated")
	}
	return nil
}
