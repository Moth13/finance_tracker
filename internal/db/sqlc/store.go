package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Store provides all function to execute db queries and transactions
type Store interface {
	Querier
	AddLineTx(ctx context.Context, arg AddLineTxParams) (AddLineTxResult, error)
	DeleteLineTx(ctx context.Context, arg DeleteLineTxParams) (DeleteLineTxResult, error)
	UpdateLineTx(ctx context.Context, arg UpdateLineTxParams) (UpdateLineTxResult, error)
}

// Store provides all function to execute db queries and transactions
type SQLStore struct {
	*Queries
	connPool *pgxpool.Pool
}

// NewStore creates a new store
func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}

// execTex executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.connPool.Begin(ctx)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
