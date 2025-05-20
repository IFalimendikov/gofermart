package storage

import (
	"context"
	"database/sql"

	"gofermart/internal/models"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func (s *Storage) GetOrders(ctx context.Context, login string) ([]models.Order, error) {
	orders := make([]models.Order, 0)
	query := `SELECT number, status, accrual, uploaded_at FROM orders WHERE login = $1 ORDER BY uploaded_at DESC`
	rows, err := s.DB.Query(query, login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var o models.Order
		var accrual sql.NullFloat64
		err := rows.Scan(&o.Order, &o.Status, &accrual, &o.UploadedAt)
		if err != nil {
			return nil, err
		}
		if accrual.Valid {
			o.Accrual = accrual.Float64
		} else {
			o.Accrual = 0.0
		}
		orders = append(orders, o)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, ErrNoOrdersFound
	}

	return orders, nil
}
