CREATE TABLE state (
    drone_id TEXT NOT NULL,
    time TIMESTAMPTZ NOT NULL,
    status TEXT NOT NULL,
    battery_pct DOUBLE PRECISION NOT NULL,
    voltage DOUBLE PRECISION NOT NULL,
    connected BOOLEAN NOT NULL,
    flight_mode TEXT NOT NULL
);

SELECT create_hypertable('state', 'time');
