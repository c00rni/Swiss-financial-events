-- name: CreateEvent :one
INSERT INTO events (id, title, link, address, location, start_at, end_at, description)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetEvents :many
SELECT *
FROM events
WHERE end_at > DATE('now')
ORDER BY start_at;

-- name: GetEventsByLink :one
SELECT *
FROM events
WHERE link = ?;

-- name: GetFilteredEvents :many
SELECT *
FROM events as e
WHERE e.location LIKE ?
AND e.end_at > DATE('now')
AND e.id IN (SELECT DISTINCT ec.event_id
                FROM eventCategories as ec, categories as c
                WHERE ec.category_id = c.id
                AND c.name LIKE ?)
AND e.id IN (SELECT DISTINCT et.event_id
            FROM eventTopics as et, topics as t
            WHERE et.topic_id = t.id
            AND t.name LIKE ?)
ORDER BY start_at;

