package storage

import (
	"context"

	"gofermart/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func (s *Storage) GetOrders(ctx context.Context, userID string) ([]models.Order, error) {
	orders := make([]models.Order, 0)
	query := `SELECT order, status, accrual, uploaded_at FROM orders WHERE user_id = $1 ORDER BY uploaded_at DESC`
	rows, err := s.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.Order, &order.Status, &order.Accrual, &order.UploadedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	if len(orders) == 0 {
		return nil, ErrNoOrdersFound
	}
	return orders, nil
}
