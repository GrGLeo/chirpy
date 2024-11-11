-- name: WriteRefreshToken :exec
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES (
  $1,
  NOW(),
  NOW(),
  $2,
  $3
);

-- name: GetRefreshToken :one
SELECT token, expires_at, revoked_at, user_id
FROM refresh_tokens
WHERE token = $1;

-- name: RevokeToken :exec
UPDATE refresh_tokens
SET updated_at = $1, revoked_at = $1
WHERE token = $2;
