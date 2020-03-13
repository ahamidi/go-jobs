package jobs

import "github.com/jackc/pgx"

var NoPendingJobsError = pgx.ErrNoRows
