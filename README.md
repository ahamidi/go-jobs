## go-jobs - Postgres Backed Jobs Queue

`go-jobs` leverages the `SKIP LOCKED` feature added in PostgreSQL 9.5 to provide a reliable and performant background job queueing system in Go.

It's built using the [pgx](https://github.com/jackc/pgx) Postgres client and requires a Postgres Database.


### Usage

See [example](/example) for a full implementation example.

Get the package:
```
go get github.com/ahamidi/go-jobs
```

```go
q, err := jobs.NewQueue("<queue_name>", "<database_connection_string>")

```

You can now create a `WorkerPool` to process that `Queue` with the following:
```go
p, err := jobs.NewWorkerPool(q, workerFunc, poolSize)
```

* `q` is the `Queue` created above
* `workerFunc` is the function to process each `Job`
* `poolSize` is an int for the number of Workers to provision in the Worker Pool.


### Benchmarks

### Contribute
