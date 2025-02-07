// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: user.sql

package db

import (
	"context"
	"time"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
  username,
  hashed_password,
  full_name,
  email,
  currency,
  password_changed_at
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING username, hashed_password, full_name, email, currency, password_changed_at, create_at
`

type CreateUserParams struct {
	Username          string    `json:"username"`
	HashedPassword    string    `json:"hashed_password"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	Currency          string    `json:"currency"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.Username,
		arg.HashedPassword,
		arg.FullName,
		arg.Email,
		arg.Currency,
		arg.PasswordChangedAt,
	)
	var i User
	err := row.Scan(
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.Currency,
		&i.PasswordChangedAt,
		&i.CreateAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users WHERE username = $1
`

func (q *Queries) DeleteUser(ctx context.Context, username string) error {
	_, err := q.db.Exec(ctx, deleteUser, username)
	return err
}

const getUser = `-- name: GetUser :one
SELECT username, hashed_password, full_name, email, currency, password_changed_at, create_at FROM users
WHERE username = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRow(ctx, getUser, username)
	var i User
	err := row.Scan(
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.Currency,
		&i.PasswordChangedAt,
		&i.CreateAt,
	)
	return i, err
}
