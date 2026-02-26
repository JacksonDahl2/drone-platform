-- name: InsertEvent :exec
INSERT INTO events (
    drone_id,
    time,
    event_type,
    payload
) VALUES (
    $1, $2, $3, $4
);

-- name: GetLatestEventsByDrone :one
SELECT * FROM events
WHERE drone_id = $1
ORDER BY time DESC
LIMIT 1;