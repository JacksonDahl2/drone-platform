-- name: InsertEvent :exec
INSERT INTO events (
    drone_id,
    time,
    event_type,
    payload
) VALUES (
    $1, $2, $3, $4
);
