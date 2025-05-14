package storage

import (
	"context"
	"fmt"
	"database/sql"
	"gofermart/internal/config"
	"gofermart/internal/models"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	cfg *config.Config
	DB  *sql.DB
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

	var usersQuery = `CREATE TABLE IF NOT EXISTS users (login text PRIMARY KEY, password text);`
	var ordersQuery = `CREATE TABLE IF NOT EXISTS orders (number text PRIMARY KEY, login text, status text, accrual FLOAT, uploaded_at text);`
	var withdrawalsQuery = `CREATE TABLE IF NOT EXISTS withdrawals (number text PRIMARY KEY, login text, sum FLOAT, processed_at text);`
	var balancesQuery = `CREATE TABLE IF NOT EXISTS balances (login text PRIMARY KEY, current FLOAT, withdrawn FLOAT);`

	tables := []string{usersQuery, ordersQuery, withdrawalsQuery, balancesQuery}

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
		var order models.Order
		err = rows.Scan(&order.Order, &order.ID, &order.Status)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
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

	var queryBal = `UPDATE balances SET current = current + $1 WHERE login = $2`
	stmtBal, err := tx.PrepareContext(ctx, queryBal)
	if err != nil {
		return err
	}
	defer stmtBal.Close()

	for _, order := range orders {
		_, err := stmtOrdr.ExecContext(ctx, order.Status, order.Accrual, order.Order)
		if err != nil {
			return err
		}
		if order.Accrual != 0 {
					fmt.Println("add accrual")
					fmt.Println(order.Accrual)
					fmt.Println(order.ID)
					fmt.Println(order.Order)
			_, err = stmtBal.ExecContext(ctx, order.Accrual, order.ID)
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

// func (s *Storage) Drop(ctx context.Context) error {
//     tx, err := s.DB.Begin()
//     if err != nil {
//         return err
//     }
//     defer tx.Rollback()

//     tables := []string{
//         "DROP TABLE IF EXISTS withdrawals",
//         "DROP TABLE IF EXISTS orders",
//         "DROP TABLE IF EXISTS balances",
//         "DROP TABLE IF EXISTS users",
//     }

//     for _, query := range tables {
//         _, err = tx.ExecContext(ctx, query)
//         if err != nil {
//             return err
//         }
//     }

//     return tx.Commit()
// }
