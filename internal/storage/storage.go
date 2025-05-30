package storage

import (
	"context"
	"database/sql"
	"gofermart/internal/config"
	"gofermart/internal/models"

	sq "github.com/Masterminds/squirrel"
)

type Storage struct {
	cfg *config.Config
	DB  *sql.DB
}

var (
	UsersQuery       = `CREATE TABLE IF NOT EXISTS users (login text PRIMARY KEY, password text);`
	OrdersQuery      = `CREATE TABLE IF NOT EXISTS orders (number text PRIMARY KEY, login text, status text, accrual FLOAT DEFAULT 0.0 NOT NULL, uploaded_at text);`
	WithdrawalsQuery = `CREATE TABLE IF NOT EXISTS withdrawals ("order" text PRIMARY KEY, login text, sum FLOAT DEFAULT 0.0 NOT NULL, processed_at text);`
	BalancesQuery    = `CREATE TABLE IF NOT EXISTS balances (login text PRIMARY KEY, current FLOAT DEFAULT 0.0 NOT NULL, withdrawn FLOAT DEFAULT 0.0 NOT NULL);`
)

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

	tables := []string{UsersQuery, OrdersQuery, WithdrawalsQuery, BalancesQuery}

	for _, q := range tables {
		_, err = db.ExecContext(ctx, q)
		if err != nil {
			return nil, err
		}
	}

	storage := Storage{
		cfg: cfg,
		DB:  db,
	}

	return &storage, nil
}

func (s *Storage) GetOrdersNums(ctx context.Context) ([]models.Order, error) {
	orders := make([]models.Order, 0)

	rows, err := sq.Select("number", "login", "status").
		From("orders").
		Where(sq.Or{
			sq.Eq{"status": "NEW"},
			sq.Eq{"status": "PROCESSING"},
		}).
		RunWith(s.DB).
		PlaceholderFormat(sq.Dollar).
		QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var o models.Order
		err = rows.Scan(&o.Order, &o.ID, &o.Status)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *Storage) UpdateOrders(ctx context.Context,  runner sq.BaseRunner, orders []models.Order) error {
	for _, o := range orders {
		_, err := sq.Update("orders").
			Set("status", o.Status).
			Set("accrual", o.Accrual).
			Where(sq.Eq{"number": o.Order}).
			RunWith(runner).
			PlaceholderFormat(sq.Dollar).
			ExecContext(ctx)
		if err != nil {
			return err
		}

		if o.Accrual != 0 {
			_, err = sq.Update("balances").
				Set("current", o.Accrual).
				Where(sq.Eq{"login": o.ID}).
				RunWith(runner).
				PlaceholderFormat(sq.Dollar).
				ExecContext(ctx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
