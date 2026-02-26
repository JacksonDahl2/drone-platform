-- name: InsertState :exec
INSERT INTO state (
    drone_id,
    time,
    status,
    battery_pct,
    voltage,
    connected,
    flight_mode
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
);

-- name: GetLatestStateByDrone :one
SELECT * FROM state
WHERE drone_id = $1
ORDER BY time DESC
LIMIT 1;