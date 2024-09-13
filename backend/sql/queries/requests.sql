-- name: CreateRequest :one
INSERT INTO requests (id, user_id, date)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetUserRequests :many
SELECT *
FROM requests
WHERE user_id = ?
AND date >= ?
ORDER BY date;

