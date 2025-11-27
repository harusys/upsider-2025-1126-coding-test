-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUsersByCompanyID :many
SELECT * FROM users WHERE company_id = $1 ORDER BY id;

-- name: CreateUser :one
INSERT INTO users (
    company_id,
    name,
    email,
    password_hash
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: UpdateUser :one
UPDATE users SET
    name = $2,
    email = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpdateUserPassword :exec
UPDATE users SET
    password_hash = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: ExistsUserByEmail :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);
