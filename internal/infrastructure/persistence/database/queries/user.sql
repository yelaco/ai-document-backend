-- name: CreateUser :one
INSERT INTO users (
    email,
    password_hash,
    full_name,
	role
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetUser :one
SELECT id, email, password_hash, full_name, role, created_at, updated_at
FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT id, email, password_hash, full_name, role, created_at, updated_at
FROM users
WHERE id = $1 LIMIT 1;
