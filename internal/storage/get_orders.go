package storage

import (
	"context"
	"gofermart/internal/models"

	sq "github.com/Masterminds/squirrel"
)

func (s *Storage) GetOrders(ctx context.Context, login string) ([]models.Order, error) {
	orders := make([]models.Order, 0)

	rows, err := sq.Select("number", "status", "accrual", "uploaded_at").
		From("orders").
		Where(sq.Eq{"login": login}).
		OrderBy("uploaded_at DESC").
		RunWith(s.DB).
		PlaceholderFormat(sq.Dollar).
		QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var o models.Order
		err := rows.Scan(&o.Order, &o.Status, &o.Accrual, &o.UploadedAt)
		if err != nil {
			return nil, err
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
