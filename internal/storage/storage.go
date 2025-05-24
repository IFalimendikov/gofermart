package storage

import (
	"context"
	"database/sql"
	"gofermart/internal/config"
	"gofermart/internal/models"
)

type Storage struct {
	cfg *config.Config
	DB  *sql.DB
}

var (
	UsersQuery = `CREATE TABLE IF NOT EXISTS users (login text PRIMARY KEY, password text);`
	OrdersQuery = `CREATE TABLE IF NOT EXISTS orders (number text PRIMARY KEY, login text, status text, accrual FLOAT DEFAULT 0, uploaded_at text);`
	WithdrawalsQuery = `CREATE TABLE IF NOT EXISTS withdrawals ("order" text PRIMARY KEY, login text, sum FLOAT DEFAULT 0, processed_at text);`
	BalancesQuery = `CREATE TABLE IF NOT EXISTS balances (login text PRIMARY KEY, current FLOAT DEFAULT 0, withdrawn FLOAT DEFAULT 0);`
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
	var query = `SELECT number, login, status FROM orders WHERE status = $1 OR status = $2`
	stmt, err := s.DB.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, "NEW", "PROCESSING")
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

func (s *Storage) UpdateOrders(ctx context.Context, orders []models.Order) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var queryOrdr = `UPDATE orders SET status = $1, accrual = $2 WHERE number = $3`
	stmtOrdr, err := tx.PrepareContext(ctx, queryOrdr)
	if err != nil {
		return err
	}
	defer stmtOrdr.Close()

	var queryBal = `UPDATE balances SET current = $1 WHERE login = $2`
	stmtBal, err := tx.PrepareContext(ctx, queryBal)
	if err != nil {
		return err
	}
	defer stmtBal.Close()

	for _, o := range orders {
		_, err := stmtOrdr.ExecContext(ctx, o.Status, o.Accrual, o.Order)
		if err != nil {
			return err
		}
		if o.Accrual != 0 {
			_, err = stmtBal.ExecContext(ctx, o.Accrual, o.ID)
			if err != nil {
				return err
			}
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}