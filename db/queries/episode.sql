-- name: FindAllEpisodes :many
SELECT * FROM episodes ORDER BY created_at DESC;

-- name: FindEpisodeByID :one
SELECT * FROM episodes WHERE id = $1;

-- name: CreateEpisode :one
INSERT INTO episodes (id, created_at, updated_at, feed_id, feed_guid, current_position, played)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateEpisode :exec
UPDATE episodes 
SET updated_at = $2, feed_id = $3, feed_guid = $4, current_position = $5, played = $6
WHERE id = $1;

-- name: DeleteEpisode :exec
DELETE FROM episodes WHERE id = $1;
