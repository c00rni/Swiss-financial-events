// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"time"
)

type Event struct {
	ID          string
	Title       string
	Link        string
	Address     string
	Location    string
	StartAt     time.Time
	EndAt       time.Time
	Description string
}
