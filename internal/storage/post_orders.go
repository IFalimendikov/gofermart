package storage

import (
	"context"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func (s *Storage) PostOrders(ctx context.Context, userID, orderNum string) error {
	var sUser string
	var sNumber string
	query := `SELECT user_id, order_id FROM orders WHERE order_id = $1`
	row := s.DB.QueryRowContext(ctx, query, orderNum)

	err := row.Scan(&sUser, &sNumber)
	if err != nil {
		return err
	}
	switch {
	case userID == sUser && orderNum == sNumber:
		return ErrDuplicateOrder
	case orderNum == sNumber && userID != sUser:
		return ErrDuplicateNumber
	}

	query = `INSERT into orders (order_id, user_id, status, uploaded_at) VALUES ($1, $2, $3, $4)`
	_, err = s.DB.ExecContext(ctx, query, orderNum, userID, "NEW", time.Now().Format(time.RFC3339))
	if err != nil {
		return err
	}
	return nil
}
