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

-- name: GetAllFilesWithPaginationAndFilters :many
SELECT * FROM files
WHERE 
    (sqlc.narg(file_name_filter)::text IS NULL OR file_name ILIKE '%' || sqlc.narg(file_name_filter)::text || '%')
    AND (sqlc.narg(mime_type_filter)::text IS NULL OR mime_type = sqlc.narg(mime_type_filter)::text)
    AND (sqlc.narg(category_filter)::text IS NULL OR category = sqlc.narg(category_filter)::text)
    AND (sqlc.narg(uploaded_by_filter)::integer IS NULL OR uploaded_by = sqlc.narg(uploaded_by_filter)::integer)
    AND (sqlc.narg(created_after)::timestamp IS NULL OR created_at >= sqlc.narg(created_after)::timestamp)
    AND (sqlc.narg(created_before)::timestamp IS NULL OR created_at <= sqlc.narg(created_before)::timestamp)
    AND (
        sqlc.narg(search)::text IS NULL 
        OR file_name ILIKE '%' || sqlc.narg(search)::text || '%' 
        OR original_name ILIKE '%' || sqlc.narg(search)::text || '%'
        OR description ILIKE '%' || sqlc.narg(search)::text || '%'
    )
ORDER BY
    CASE WHEN @sort_field::text = 'id' AND @sort_order::text = 'ASC' THEN id END ASC,
    CASE WHEN @sort_field::text = 'id' AND @sort_order::text = 'DESC' THEN id END DESC,
    CASE WHEN @sort_field::text = 'file_name' AND @sort_order::text = 'ASC' THEN file_name END ASC,
    CASE WHEN @sort_field::text = 'file_name' AND @sort_order::text = 'DESC' THEN file_name END DESC,
    CASE WHEN @sort_field::text = 'file_size' AND @sort_order::text = 'ASC' THEN file_size END ASC,
    CASE WHEN @sort_field::text = 'file_size' AND @sort_order::text = 'DESC' THEN file_size END DESC,
    CASE WHEN @sort_field::text = 'created_at' AND @sort_order::text = 'ASC' THEN created_at END ASC,
    CASE WHEN @sort_field::text = 'created_at' AND @sort_order::text = 'DESC' THEN created_at END DESC,
    created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetFilesByUserWithPagination :many
SELECT * FROM files
WHERE uploaded_by = $3
    AND (sqlc.narg(file_name_filter)::text IS NULL OR file_name ILIKE '%' || sqlc.narg(file_name_filter)::text || '%')
    AND (sqlc.narg(mime_type_filter)::text IS NULL OR mime_type = sqlc.narg(mime_type_filter)::text)
    AND (sqlc.narg(category_filter)::text IS NULL OR category = sqlc.narg(category_filter)::text)
    AND (sqlc.narg(created_after)::timestamp IS NULL OR created_at >= sqlc.narg(created_after)::timestamp)
    AND (sqlc.narg(created_before)::timestamp IS NULL OR created_at <= sqlc.narg(created_before)::timestamp)
    AND (
        sqlc.narg(search)::text IS NULL 
        OR file_name ILIKE '%' || sqlc.narg(search)::text || '%' 
        OR original_name ILIKE '%' || sqlc.narg(search)::text || '%'
        OR description ILIKE '%' || sqlc.narg(search)::text || '%'
    )
ORDER BY
    CASE WHEN @sort_field::text = 'id' AND @sort_order::text = 'ASC' THEN id END ASC,
    CASE WHEN @sort_field::text = 'id' AND @sort_order::text = 'DESC' THEN id END DESC,
    CASE WHEN @sort_field::text = 'file_name' AND @sort_order::text = 'ASC' THEN file_name END ASC,
    CASE WHEN @sort_field::text = 'file_name' AND @sort_order::text = 'DESC' THEN file_name END DESC,
    CASE WHEN @sort_field::text = 'file_size' AND @sort_order::text = 'ASC' THEN file_size END ASC,
    CASE WHEN @sort_field::text = 'file_size' AND @sort_order::text = 'DESC' THEN file_size END DESC,
    CASE WHEN @sort_field::text = 'created_at' AND @sort_order::text = 'ASC' THEN created_at END ASC,
    CASE WHEN @sort_field::text = 'created_at' AND @sort_order::text = 'DESC' THEN created_at END DESC,
    created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountFiles :one
SELECT COUNT(*) FROM files;

-- name: CountFilesWithFilters :one
SELECT COUNT(*) FROM files
WHERE 
    (sqlc.narg(file_name_filter)::text IS NULL OR file_name ILIKE '%' || sqlc.narg(file_name_filter)::text || '%')
    AND (sqlc.narg(mime_type_filter)::text IS NULL OR mime_type = sqlc.narg(mime_type_filter)::text)
    AND (sqlc.narg(category_filter)::text IS NULL OR category = sqlc.narg(category_filter)::text)
    AND (sqlc.narg(uploaded_by_filter)::integer IS NULL OR uploaded_by = sqlc.narg(uploaded_by_filter)::integer)
    AND (sqlc.narg(created_after)::timestamp IS NULL OR created_at >= sqlc.narg(created_after)::timestamp)
    AND (sqlc.narg(created_before)::timestamp IS NULL OR created_at <= sqlc.narg(created_before)::timestamp)
    AND (
        sqlc.narg(search)::text IS NULL 
        OR file_name ILIKE '%' || sqlc.narg(search)::text || '%' 
        OR original_name ILIKE '%' || sqlc.narg(search)::text || '%'
        OR description ILIKE '%' || sqlc.narg(search)::text || '%'
    );

-- name: CountFilesByUser :one
SELECT COUNT(*) FROM files
WHERE uploaded_by = $1
    AND (sqlc.narg(file_name_filter)::text IS NULL OR file_name ILIKE '%' || sqlc.narg(file_name_filter)::text || '%')
    AND (sqlc.narg(mime_type_filter)::text IS NULL OR mime_type = sqlc.narg(mime_type_filter)::text)
    AND (sqlc.narg(category_filter)::text IS NULL OR category = sqlc.narg(category_filter)::text)
    AND (sqlc.narg(created_after)::timestamp IS NULL OR created_at >= sqlc.narg(created_after)::timestamp)
    AND (sqlc.narg(created_before)::timestamp IS NULL OR created_at <= sqlc.narg(created_before)::timestamp)
    AND (
        sqlc.narg(search)::text IS NULL 
        OR file_name ILIKE '%' || sqlc.narg(search)::text || '%' 
        OR original_name ILIKE '%' || sqlc.narg(search)::text || '%'
        OR description ILIKE '%' || sqlc.narg(search)::text || '%'
    );

-- name: UpdateFile :one
UPDATE files
SET description = $2, category = $3, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteFile :exec
DELETE FROM files
WHERE id = $1;