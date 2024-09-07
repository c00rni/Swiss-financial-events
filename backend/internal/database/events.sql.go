// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: events.sql

package database

import (
	"context"
	"time"
)

const createEvent = `-- name: CreateEvent :one
INSERT INTO events (id, title, link, address, location, start_at, end_at, description)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING id, title, link, address, location, start_at, end_at, description
`

type CreateEventParams struct {
	ID          string
	Title       string
	Link        string
	Address     string
	Location    string
	StartAt     time.Time
	EndAt       time.Time
	Description string
}

func (q *Queries) CreateEvent(ctx context.Context, arg CreateEventParams) (Event, error) {
	row := q.db.QueryRowContext(ctx, createEvent,
		arg.ID,
		arg.Title,
		arg.Link,
		arg.Address,
		arg.Location,
		arg.StartAt,
		arg.EndAt,
		arg.Description,
	)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Link,
		&i.Address,
		&i.Location,
		&i.StartAt,
		&i.EndAt,
		&i.Description,
	)
	return i, err
}

const getEvents = `-- name: GetEvents :many
SELECT id, title, link, address, location, start_at, end_at, description
FROM events
ORDER BY start_at
`

func (q *Queries) GetEvents(ctx context.Context) ([]Event, error) {
	rows, err := q.db.QueryContext(ctx, getEvents)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Event
	for rows.Next() {
		var i Event
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Link,
			&i.Address,
			&i.Location,
			&i.StartAt,
			&i.EndAt,
			&i.Description,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
