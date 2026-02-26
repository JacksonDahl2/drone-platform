-- name: GetLatestGpsAllDrones :many
SELECT DISTINCT ON (drone_id) *
FROM gps
ORDER BY drone_id, time DESC;

-- name: GetLatestStateAllDrones :many
SELECT DISTINCT ON (drone_id) *
FROM state
ORDER BY drone_id, time DESC;

-- name: GetLatestGpsAndStatePerDrone :many
WITH latest_gps AS (
    SELECT DISTINCT ON (drone_id)
        drone_id, time, latitude, longitude, altitude, heading, pitch, roll, speed, climb_rate, angular_rate
    FROM gps
    ORDER BY drone_id, time DESC
),
latest_state AS (
    SELECT DISTINCT ON (drone_id)
        drone_id, time, status, battery_pct, voltage, connected, flight_mode
    FROM state
    ORDER BY drone_id, time DESC
)
SELECT
    g.drone_id,
    g.time AS gps_time,
    g.latitude,
    g.longitude,
    g.altitude,
    g.heading,
    g.pitch,
    g.roll,
    g.speed,
    g.climb_rate,
    g.angular_rate,
    s.time AS state_time,
    s.status,
    s.battery_pct,
    s.voltage,
    s.connected,
    s.flight_mode
FROM latest_gps g
JOIN latest_state s ON g.drone_id = s.drone_id;

-- name: GetGpsByDroneTimeRange :many
SELECT * FROM gps
WHERE drone_id = $1 AND time >= $2 AND time <= $3
ORDER BY time ASC;

-- name: GetStateByDroneTimeRange :many
SELECT * FROM state
WHERE drone_id = $1 AND time >= $2 AND time <= $3
ORDER BY time ASC;

-- name: GetEventsByTimeRange :many
SELECT * FROM events
WHERE time >= $1 AND time <= $2
ORDER BY time DESC;

-- name: GetEventsByDroneTimeRange :many
SELECT * FROM events
WHERE drone_id = $1 AND time >= $2 AND time <= $3
ORDER BY time DESC;

-- name: GetRecentEvents :many
SELECT * FROM events
ORDER BY time DESC
LIMIT $1;

-- name: GetDroneCount :one
SELECT COUNT(DISTINCT drone_id)::bigint AS count FROM gps;

-- name: GetDroneCountByStatus :many
WITH latest AS (
    SELECT DISTINCT ON (drone_id) drone_id, status
    FROM state
    ORDER BY drone_id, time DESC
)
SELECT status, COUNT(*)::bigint AS count
FROM latest
GROUP BY status;

-- name: GetConnectedDroneCount :one
WITH latest AS (
    SELECT DISTINCT ON (drone_id) drone_id, connected
    FROM state
    ORDER BY drone_id, time DESC
)
SELECT COUNT(*)::bigint AS count FROM latest WHERE connected = true;

-- name: GetEventCountByType :many
SELECT event_type, COUNT(*)::bigint AS count
FROM events
WHERE time >= $1 AND time <= $2
GROUP BY event_type
ORDER BY count DESC;

-- name: GetBatteryStatsByDrone :many
SELECT
    drone_id,
    MIN(battery_pct) AS min_pct,
    AVG(battery_pct) AS avg_pct,
    MAX(battery_pct) AS max_pct
FROM state
WHERE time >= $1 AND time <= $2
GROUP BY drone_id;

-- name: GetActivityByTimeBucket :many
SELECT
    time_bucket($1::interval, time) AS bucket,
    COUNT(*)::bigint AS count
FROM gps
WHERE time >= $2 AND time <= $3
GROUP BY bucket
ORDER BY bucket;

-- name: GetActiveDronesPerTimeBucket :many
SELECT
    time_bucket($1::interval, time) AS bucket,
    COUNT(DISTINCT drone_id)::bigint AS drone_count
FROM gps
WHERE time >= $2 AND time <= $3
GROUP BY bucket
ORDER BY bucket;
