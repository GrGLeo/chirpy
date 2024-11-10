-- name: GetPassword :one
SELECT hashed_password
FROM users
WHERE id = $1;
