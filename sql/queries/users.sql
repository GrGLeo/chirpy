-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
  gen_random_uuid(),
  NOW(),
  NOW(),
  $1,
  $2
)
RETURNING id, created_at, updated_at, email;

-- name: GetHashedPassword :one
SELECT *
FROM users
WHERE email = $1;

-- name: UpdateUserInfo :one
UPDATE users
SET hashed_password = $1, email = $2, updated_at = NOW()
WHERE id = $3
RETURNING id, created_at, updated_at, email;

-- name: GetRed :one
SELECT is_chirpy_red
FROM users
WHERE id = $1;

-- name: UpgradeUser :exec
UPDATE users
SET is_chirpy_red = true
WHERE id = $1;
