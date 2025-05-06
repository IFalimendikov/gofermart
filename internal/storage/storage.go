package storage

import (
	// "bufio"
	"context"
	// "encoding/json"
	// "errors"
	// "os"

	"database/sql"
	"gofermart/internal/config"
	// "gofermart/internal/models"

	// "github.com/deatil/go-encoding/base62"
	// "github.com/jackc/pgerrcode"
	// "github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	cfg *config.Config
	DB *sql.DB
}

func New(ctx context.Context, cfg *config.Config) (*Storage, error) {
	if cfg.DatabaseURI == "" {
		return nil, ErrBadConn
	}

	db, err := sql.Open("pgx", cfg.DatabaseURI)
	if err != nil {
		return nil, ErrBadConn
	}

	err = db.Ping()
	if err != nil {
		return nil, ErrBadConn
	}

	var usersQuery = `CREATE TABLE IF NOT EXISTS users (user_id text PRIMARY KEY, login text UNIQUE, password text, connected bool DEFAULT false);`
	var ordersQuery = `CREATE TABLE IF NOT EXISTS orders (order_id text PRIMARY KEY);`
	var txsQuery = `CREATE TABLE IF NOT EXISTS txs (tx_id text PRIMARY KEY);`
	var balancesQuery = `CREATE TABLE IF NOT EXISTS balances (balance_id text PRIMARY KEY);`

	tables := []string{usersQuery, ordersQuery, txsQuery, balancesQuery}

	for _, q := range tables {
		_, err = db.ExecContext(ctx, q)
		if err != nil {
			return nil, err
		}
	}

	storage := Storage{
		cfg: cfg,
		DB: db,
	}

	return &storage, nil
}