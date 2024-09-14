-- name: LinkEventToTopic :one
INSERT INTO eventTopics (event_id, topic_id)
VALUES (?, ?)
RETURNING *;

