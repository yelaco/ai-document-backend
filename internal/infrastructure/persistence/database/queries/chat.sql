-- name: CreateChat :one
INSERT INTO chats (

) VALUES (
) RETURNING *;

-- name: GetPaginatedChatsByUserID :many
SELECT *
FROM chats
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetPaginatedChatsByDocumentID :many
SELECT *
FROM chats
WHERE document_id = $1 AND user_id = $2
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: CountUserChats :one
SELECT COUNT(*)
FROM chats
WHERE user_id = $1;

-- name: CountDocumentChats :one
SELECT COUNT(*)
FROM chats
WHERE document_id = $1 AND user_id = $1;

-- name: GetChat :one
SELECT *
FROM chats
WHERE id = $1 AND user_id = $2
LIMIT 1;

-- name: DeleteChat :exec
DELETE FROM chats
WHERE id = $1 AND user_id = $2;
