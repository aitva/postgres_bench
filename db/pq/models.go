// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0

package pq

import (
	"time"

	"github.com/google/uuid"
)

type Page struct {
	ID        uuid.UUID
	UpdatedAt time.Time
	Title     string
	Text      string
}