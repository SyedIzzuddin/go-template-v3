-- name: CreateFile :one
INSERT INTO files (file_name, original_name, file_path, file_size, mime_type, description, category, uploaded_by)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetFile :one
SELECT * FROM files
WHERE id = $1 LIMIT 1;

-- name: GetFilesByUser :many
SELECT * FROM files
WHERE uploaded_by = $1
ORDER BY created_at DESC;

-- name: GetAllFiles :many
SELECT * FROM files
ORDER BY created_at DESC;

-- name: UpdateFile :one
UPDATE files
SET description = $2, category = $3, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteFile :exec
DELETE FROM files
WHERE id = $1;