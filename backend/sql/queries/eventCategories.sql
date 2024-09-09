-- name: LinkEventToCategory :one
INSERT INTO eventCategories (event_id, category_id)
VALUES (?, ?)
RETURNING *;

