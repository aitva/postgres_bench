// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: query.sql

package pgx

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createPage = `-- name: CreatePage :exec
INSERT INTO pages (id, updated_at, title, text)
VALUES ($1, $2, $3, $4)
`

type CreatePageParams struct {
	ID        uuid.UUID
	UpdatedAt time.Time
	Title     string
	Text      string
}

func (q *Queries) CreatePage(ctx context.Context, arg CreatePageParams) error {
	_, err := q.db.Exec(ctx, createPage,
		arg.ID,
		arg.UpdatedAt,
		arg.Title,
		arg.Text,
	)
	return err
}

const getPage = `-- name: GetPage :one
SELECT id, updated_at, title, text FROM pages WHERE id = $1
`

func (q *Queries) GetPage(ctx context.Context, id uuid.UUID) (Page, error) {
	row := q.db.QueryRow(ctx, getPage, id)
	var i Page
	err := row.Scan(
		&i.ID,
		&i.UpdatedAt,
		&i.Title,
		&i.Text,
	)
	return i, err
}

const listPageIDs = `-- name: ListPageIDs :many
SELECT id FROM pages
WHERE $1::uuid IS NULL OR id > $1
LIMIT $2
`

type ListPageIDsParams struct {
	Cursor uuid.NullUUID
	Limit  int32
}

func (q *Queries) ListPageIDs(ctx context.Context, arg ListPageIDsParams) ([]uuid.UUID, error) {
	rows, err := q.db.Query(ctx, listPageIDs, arg.Cursor, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listPages = `-- name: ListPages :many
SELECT id, updated_at, title, text FROM pages
WHERE $1::uuid IS NULL or id > $1
LIMIT $2
`

type ListPagesParams struct {
	Cursor uuid.NullUUID
	Limit  int32
}

func (q *Queries) ListPages(ctx context.Context, arg ListPagesParams) ([]Page, error) {
	rows, err := q.db.Query(ctx, listPages, arg.Cursor, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Page
	for rows.Next() {
		var i Page
		if err := rows.Scan(
			&i.ID,
			&i.UpdatedAt,
			&i.Title,
			&i.Text,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
