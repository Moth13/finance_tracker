// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: recline.sql

package db

import (
	"context"
	"time"

	decimal "github.com/shopspring/decimal"
)

const createRecLine = `-- name: CreateRecLine :one
INSERT INTO reclines (
  title,
  owner,
  account_id,
  category_id,
  amount,
  description,
  recurrency,
  due_date
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING id, owner, title, account_id, amount, category_id, description, recurrency, due_date
`

type CreateRecLineParams struct {
	Title       string          `json:"title"`
	Owner       string          `json:"owner"`
	AccountID   int64           `json:"account_id"`
	CategoryID  int64           `json:"category_id"`
	Amount      decimal.Decimal `json:"amount"`
	Description string          `json:"description"`
	Recurrency  string          `json:"recurrency"`
	DueDate     time.Time       `json:"due_date"`
}

func (q *Queries) CreateRecLine(ctx context.Context, arg CreateRecLineParams) (Recline, error) {
	row := q.db.QueryRow(ctx, createRecLine,
		arg.Title,
		arg.Owner,
		arg.AccountID,
		arg.CategoryID,
		arg.Amount,
		arg.Description,
		arg.Recurrency,
		arg.DueDate,
	)
	var i Recline
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Title,
		&i.AccountID,
		&i.Amount,
		&i.CategoryID,
		&i.Description,
		&i.Recurrency,
		&i.DueDate,
	)
	return i, err
}

const deleteRecLine = `-- name: DeleteRecLine :exec
DELETE FROM reclines WHERE id = $1
`

func (q *Queries) DeleteRecLine(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteRecLine, id)
	return err
}

const getRecLine = `-- name: GetRecLine :one
SELECT id, owner, title, account_id, amount, category_id, description, recurrency, due_date FROM reclines
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetRecLine(ctx context.Context, id int64) (Recline, error) {
	row := q.db.QueryRow(ctx, getRecLine, id)
	var i Recline
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Title,
		&i.AccountID,
		&i.Amount,
		&i.CategoryID,
		&i.Description,
		&i.Recurrency,
		&i.DueDate,
	)
	return i, err
}

const getRecLineForUpdate = `-- name: GetRecLineForUpdate :one
SELECT id, owner, title, account_id, amount, category_id, description, recurrency, due_date FROM reclines
WHERE id = $1 LIMIT 1 FOR NO KEY UPDATE
`

func (q *Queries) GetRecLineForUpdate(ctx context.Context, id int64) (Recline, error) {
	row := q.db.QueryRow(ctx, getRecLineForUpdate, id)
	var i Recline
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Title,
		&i.AccountID,
		&i.Amount,
		&i.CategoryID,
		&i.Description,
		&i.Recurrency,
		&i.DueDate,
	)
	return i, err
}

const listRecLines = `-- name: ListRecLines :many
SELECT id, owner, title, account_id, amount, category_id, description, recurrency, due_date FROM reclines
WHERE owner = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListRecLinesParams struct {
	Owner  string `json:"owner"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

func (q *Queries) ListRecLines(ctx context.Context, arg ListRecLinesParams) ([]Recline, error) {
	rows, err := q.db.Query(ctx, listRecLines, arg.Owner, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Recline{}
	for rows.Next() {
		var i Recline
		if err := rows.Scan(
			&i.ID,
			&i.Owner,
			&i.Title,
			&i.AccountID,
			&i.Amount,
			&i.CategoryID,
			&i.Description,
			&i.Recurrency,
			&i.DueDate,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateRecLine = `-- name: UpdateRecLine :one
UPDATE reclines
SET title = $2, account_id = $3, category_id = $4, amount = $5, description = $6, recurrency = $7, due_date = $8
WHERE id = $1
RETURNING id, owner, title, account_id, amount, category_id, description, recurrency, due_date
`

type UpdateRecLineParams struct {
	ID          int64           `json:"id"`
	Title       string          `json:"title"`
	AccountID   int64           `json:"account_id"`
	CategoryID  int64           `json:"category_id"`
	Amount      decimal.Decimal `json:"amount"`
	Description string          `json:"description"`
	Recurrency  string          `json:"recurrency"`
	DueDate     time.Time       `json:"due_date"`
}

func (q *Queries) UpdateRecLine(ctx context.Context, arg UpdateRecLineParams) (Recline, error) {
	row := q.db.QueryRow(ctx, updateRecLine,
		arg.ID,
		arg.Title,
		arg.AccountID,
		arg.CategoryID,
		arg.Amount,
		arg.Description,
		arg.Recurrency,
		arg.DueDate,
	)
	var i Recline
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Title,
		&i.AccountID,
		&i.Amount,
		&i.CategoryID,
		&i.Description,
		&i.Recurrency,
		&i.DueDate,
	)
	return i, err
}
