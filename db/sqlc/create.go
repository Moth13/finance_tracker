package db

import (
	"context"
	"sync"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	pgx "github.com/jackc/pgx/v5"
	pgxpool "github.com/jackc/pgx/v5/pgxpool"
)

var dtMutex sync.RWMutex

// CreateDBConnection generates from the db source uri the connector
func CreateDBConnection(uri string) (*pgxpool.Pool, error) {
	// load the config from the uri
	config, err := pgxpool.ParseConfig(uri)
	if err != nil {
		return nil, err
	}

	// Register pgxdecimal type before usage but after connection
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		dtMutex.RLock()
		defer dtMutex.RUnlock()
		pgxdecimal.Register(conn.TypeMap())
		return nil
	}

	// Create connection
	conn, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
