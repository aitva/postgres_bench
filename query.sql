-- name: CreatePage :exec
INSERT INTO pages (id, updated_at, title, text)
VALUES ($1, $2, $3, $4);

-- name: GetPage :one
SELECT * FROM pages WHERE id = $1;

-- name: ListPages :many
SELECT * FROM pages
WHERE sqlc.narg('cursor')::uuid IS NULL or id > sqlc.narg('cursor')
LIMIT sqlc.arg('limit');

-- name: ListPageIDs :many
SELECT id FROM pages
WHERE sqlc.narg('cursor')::uuid IS NULL OR id > sqlc.narg('cursor')
LIMIT sqlc.arg('limit');
