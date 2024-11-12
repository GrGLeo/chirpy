// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
  gen_random_uuid(),
  NOW(),
  NOW(),
  $1,
  $2
)
RETURNING id, created_at, updated_at, email
`

type CreateUserParams struct {
	Email          string
	HashedPassword string
}

type CreateUserRow struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Email     string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Email, arg.HashedPassword)
	var i CreateUserRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
	)
	return i, err
}

const getHashedPassword = `-- name: GetHashedPassword :one
SELECT id, created_at, updated_at, email, hashed_password, is_chirpy_red
FROM users
WHERE email = $1
`

func (q *Queries) GetHashedPassword(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getHashedPassword, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}

const getRed = `-- name: GetRed :one
SELECT is_chirpy_red
FROM users
WHERE id = $1
`

func (q *Queries) GetRed(ctx context.Context, id uuid.UUID) (bool, error) {
	row := q.db.QueryRowContext(ctx, getRed, id)
	var is_chirpy_red bool
	err := row.Scan(&is_chirpy_red)
	return is_chirpy_red, err
}

const updateUserInfo = `-- name: UpdateUserInfo :one
UPDATE users
SET hashed_password = $1, email = $2, updated_at = NOW()
WHERE id = $3
RETURNING id, created_at, updated_at, email
`

type UpdateUserInfoParams struct {
	HashedPassword string
	Email          string
	ID             uuid.UUID
}

type UpdateUserInfoRow struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Email     string
}

func (q *Queries) UpdateUserInfo(ctx context.Context, arg UpdateUserInfoParams) (UpdateUserInfoRow, error) {
	row := q.db.QueryRowContext(ctx, updateUserInfo, arg.HashedPassword, arg.Email, arg.ID)
	var i UpdateUserInfoRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
	)
	return i, err
}

const upgradeUser = `-- name: UpgradeUser :exec
UPDATE users
SET is_chirpy_red = true
WHERE id = $1
`

func (q *Queries) UpgradeUser(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, upgradeUser, id)
	return err
}
