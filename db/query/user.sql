-- name: CreateUser :one
INSERT INTO users (username, hashed_password, full_name, email, role)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET hashed_password = COALESCE(sqlc.narg('hashed_password'), hashed_password),
    full_name = COALESCE(sqlc.narg('full_name'), full_name),
    email = COALESCE(sqlc.narg('email'), email),
    is_email_verified = COALESCE(sqlc.narg('is_email_verified'), is_email_verified),
    password_changed_at = COALESCE(sqlc.narg('password_changed_at'), password_changed_at)
WHERE username = sqlc.arg('username')
RETURNING *;

-- name: UpdateUserPassword :one
UPDATE users
SET hashed_password = $2,
    password_changed_at = $3
WHERE username = $1
RETURNING *;