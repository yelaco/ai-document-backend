-- name: InsertRefreshToken :exec
INSERT INTO refresh_tokens (
	user_id,
	token_hash,
	expires_at
) VALUES (
  $1, $2, $3
);

-- name: GetRefreshTokenHashByUserID :one
SELECT id, token_hash
FROM refresh_tokens
WHERE user_id = $1 AND revoked = false AND expires_at > NOW()
ORDER BY updated_at DESC
LIMIT 1;

-- name: RevokeRefreshTokenByUserID :exec
UPDATE refresh_tokens
SET revoked = TRUE, revoked_at = NOW()
WHERE user_id = $1 AND revoked = false;

-- name: DeleteRefreshTokenByUserID :exec
DELETE FROM refresh_tokens
WHERE user_id = $1;
