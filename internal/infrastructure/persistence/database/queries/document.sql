-- name: CreateDocument :one
INSERT INTO documents (
    title,
    status,
    user_id
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetPaginatedDocuments :many
SELECT id, title, user_id, status, created_at, updated_at
FROM documents
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountUserDocuments :one
SELECT COUNT(*)
FROM documents
WHERE user_id = $1;

-- name: GetDocument :one
SELECT id, title, user_id, status, created_at, updated_at
FROM documents
WHERE id = $1 AND user_id = $2
LIMIT 1;

-- name: UpdateDocumentStatus :one
UPDATE documents
SET status = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteDocument :exec
DELETE FROM documents
WHERE id = $1 AND user_id = $2;
