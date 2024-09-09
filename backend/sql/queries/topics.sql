-- name: AddTopic :one
INSERT INTO topics (id, name)
VALUES (?, ?)
RETURNING *;

-- name: GetTopicByName :one
SELECT *
FROM topics
WHERE name = ?;
