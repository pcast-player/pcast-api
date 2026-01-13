-- name: FindAllUsers :many
SELECT * FROM users ORDER BY created_at DESC;

-- name: FindUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: FindUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, password)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateUser :exec
UPDATE users
SET updated_at = $2, email = $3, password = $4
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
