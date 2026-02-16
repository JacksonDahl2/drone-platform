-- name: InsertGps :exec
INSERT INTO gps (
    drone_id,
    time,
    latitude,
    longitude,
    altitude,
    heading,
    pitch,
    roll,
    speed,
    climb_rate,
    angular_rate
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
);

-- name: GetLatestGpsByDrone :one
SELECT * FROM gps
WHERE drone_id = $1
ORDER BY time DESC
LIMIT 1;