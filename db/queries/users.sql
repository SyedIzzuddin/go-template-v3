-- name: CreateUser :one
INSERT INTO users (name, email)
VALUES ($1, $2)
RETURNING *;

-- name: CreateUserWithPassword :one
INSERT INTO users (name, email, password_hash, role, email_verification_token, email_verification_expires_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserByEmailWithPassword :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: UpdateUser :one
UPDATE users
SET name = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: GetAllUsers :many
SELECT * FROM users
ORDER BY created_at DESC;

-- name: GetUserByVerificationToken :one
SELECT * FROM users
WHERE email_verification_token = $1 LIMIT 1;

-- name: UpdateEmailVerification :one
UPDATE users
SET email_verified = $2, email_verification_token = NULL, email_verification_expires_at = NULL, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateEmailVerificationToken :one
UPDATE users
SET email_verification_token = $2, email_verification_expires_at = $3, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: VerifyEmailByToken :exec
UPDATE users
SET email_verified = true, email_verification_token = NULL, email_verification_expires_at = NULL, updated_at = NOW()
WHERE email_verification_token = $1;

-- name: UpdateVerificationToken :exec
UPDATE users
SET email_verification_token = $2, email_verification_expires_at = $3, updated_at = NOW()
WHERE id = $1;

-- name: GetUserByPasswordResetToken :one
SELECT * FROM users
WHERE password_reset_token = $1 LIMIT 1;

-- name: UpdatePasswordResetToken :exec
UPDATE users
SET password_reset_token = $2, password_reset_expires_at = $3, updated_at = NOW()
WHERE id = $1;

-- name: ResetPassword :exec
UPDATE users
SET password_hash = $2, password_reset_token = NULL, password_reset_expires_at = NULL, updated_at = NOW()
WHERE password_reset_token = $1;
