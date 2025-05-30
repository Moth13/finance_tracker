// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: account.sql

package db

import (
	"context"

	decimal "github.com/shopspring/decimal"
)

const addAccountBalance = `-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + $1, final_balance = final_balance + $2
WHERE id = $3
RETURNING id, owner, title, description, init_balance, balance, final_balance
`

type AddAccountBalanceParams struct {
	Amount      decimal.Decimal `json:"amount"`
	FinalAmount decimal.Decimal `json:"final_amount"`
	ID          int64           `json:"id"`
}

func (q *Queries) AddAccountBalance(ctx context.Context, arg AddAccountBalanceParams) (Account, error) {
	row := q.db.QueryRow(ctx, addAccountBalance, arg.Amount, arg.FinalAmount, arg.ID)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Title,
		&i.Description,
		&i.InitBalance,
		&i.Balance,
		&i.FinalBalance,
	)
	return i, err
}

const createAccount = `-- name: CreateAccount :one
INSERT INTO accounts (
  owner,
  title,
  description,
  init_balance
) VALUES (
    $1, $2, $3, $4
) RETURNING id, owner, title, description, init_balance, balance, final_balance
`

type CreateAccountParams struct {
	Owner       string          `json:"owner"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	InitBalance decimal.Decimal `json:"init_balance"`
}

func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error) {
	row := q.db.QueryRow(ctx, createAccount,
		arg.Owner,
		arg.Title,
		arg.Description,
		arg.InitBalance,
	)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Title,
		&i.Description,
		&i.InitBalance,
		&i.Balance,
		&i.FinalBalance,
	)
	return i, err
}

const deleteAccount = `-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1
`

func (q *Queries) DeleteAccount(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteAccount, id)
	return err
}

const getAccount = `-- name: GetAccount :one
SELECT id, owner, title, description, init_balance, balance, final_balance FROM accounts
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetAccount(ctx context.Context, id int64) (Account, error) {
	row := q.db.QueryRow(ctx, getAccount, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Title,
		&i.Description,
		&i.InitBalance,
		&i.Balance,
		&i.FinalBalance,
	)
	return i, err
}

const getAccountForUpdate = `-- name: GetAccountForUpdate :one
SELECT id, owner, title, description, init_balance, balance, final_balance FROM accounts
WHERE id = $1 LIMIT 1 FOR NO KEY UPDATE
`

func (q *Queries) GetAccountForUpdate(ctx context.Context, id int64) (Account, error) {
	row := q.db.QueryRow(ctx, getAccountForUpdate, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Title,
		&i.Description,
		&i.InitBalance,
		&i.Balance,
		&i.FinalBalance,
	)
	return i, err
}

const listAccounts = `-- name: ListAccounts :many
SELECT id, owner, title, description, init_balance, balance, final_balance FROM accounts
WHERE owner = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListAccountsParams struct {
	Owner  string `json:"owner"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

func (q *Queries) ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error) {
	rows, err := q.db.Query(ctx, listAccounts, arg.Owner, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Account{}
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.ID,
			&i.Owner,
			&i.Title,
			&i.Description,
			&i.InitBalance,
			&i.Balance,
			&i.FinalBalance,
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

const updateAccount = `-- name: UpdateAccount :one
UPDATE accounts
SET init_balance = $2, title = $3, description = $4
WHERE id = $1
RETURNING id, owner, title, description, init_balance, balance, final_balance
`

type UpdateAccountParams struct {
	ID          int64           `json:"id"`
	InitBalance decimal.Decimal `json:"init_balance"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
}

func (q *Queries) UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error) {
	row := q.db.QueryRow(ctx, updateAccount,
		arg.ID,
		arg.InitBalance,
		arg.Title,
		arg.Description,
	)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Title,
		&i.Description,
		&i.InitBalance,
		&i.Balance,
		&i.FinalBalance,
	)
	return i, err
}
