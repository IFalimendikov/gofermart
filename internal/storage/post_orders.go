package storage

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func (s *Storage) PostOrders(ctx context.Context, login, orderNum string) error {
	var sUser string
	var sNumber string
	query := `SELECT login, number FROM orders WHERE number = $1`
	row := s.DB.QueryRowContext(ctx, query, orderNum)

	row.Scan(&sUser, &sNumber)

	switch {
	case login == sUser && orderNum == sNumber:
		return ErrDuplicateOrder
	case orderNum == sNumber && login != sUser:
		return ErrDuplicateNumber
	}

	_, err := sq.Insert("orders").
		Columns("number", "login", "status", "uploaded_at").
		Values(orderNum, login, "NEW", time.Now().Format(time.RFC3339)).
		RunWith(s.DB).
		PlaceholderFormat(sq.Dollar).
		ExecContext(ctx)

	return err
}
