-- name: FindAllFeeds :many
SELECT * FROM feeds ORDER BY created_at DESC;

-- name: FindFeedByID :one
SELECT * FROM feeds WHERE id = $1;

-- name: FindFeedsByUserID :many
SELECT * FROM feeds WHERE user_id = $1 ORDER BY created_at DESC;

-- name: FindFeedByIDAndUserID :one
SELECT * FROM feeds WHERE id = $1 AND user_id = $2;

-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, user_id, title, url, synced_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateFeed :exec
UPDATE feeds
SET updated_at = $2, user_id = $3, title = $4, url = $5, synced_at = $6
WHERE id = $1;

-- name: DeleteFeed :exec
DELETE FROM feeds WHERE id = $1;
