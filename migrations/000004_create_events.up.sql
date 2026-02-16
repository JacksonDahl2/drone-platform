CREATE TABLE events (
    drone_id TEXT NOT NULL,
    time TIMESTAMPTZ NOT NULL,
    event_type TEXT NOT NULL,
    payload JSONB NOT NULL
);

SELECT create_hypertable('events', 'time');
