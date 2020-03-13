CREATE UNLOGGED TABLE jobs (
    id serial primary key,
    queue text NOT NULL DEFAULT 'default',
    retries int NOT NULL DEFAULT 0,
    payload jsonb NOT NULL,
    state int NOT NULL DEFAULT 0,
    success boolean,
    error text,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, 
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, 
    completed_at timestamptz 
);
create index idx_state_updated ON jobs (queue, state, updated_at ASC);
