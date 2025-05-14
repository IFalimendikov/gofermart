package storage

import (
	"context"
	"database/sql"

	"gofermart/internal/models"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func (s *Storage) GetOrders(ctx context.Context, userID string) ([]models.Order, error) {
	orders := make([]models.Order, 0)
	query := `SELECT order_id, status, accrual, uploaded_at FROM orders WHERE login = $1 ORDER BY uploaded_at DESC`
	rows, err := s.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order
		var accrual sql.NullFloat64
		err := rows.Scan(&order.Order, &order.Status, &accrual, &order.UploadedAt)
		if err != nil {
			return nil, err
		}
		if accrual.Valid {
			order.Accrual = accrual.Float64
		}
		orders = append(orders, order)
	}
	if len(orders) == 0 {
		return nil, ErrNoOrdersFound
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}
