-- name: CreateEvent :one
INSERT INTO events (id, title, link, address, location, start_at, end_at, description)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetEvents :many
SELECT *
FROM events
ORDER BY start_at;
