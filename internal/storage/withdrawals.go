package storage

import (
	"context"
	"gofermart/internal/models"

	sq "github.com/Masterminds/squirrel"
)

func (s *Storage) Withdrawals(ctx context.Context, login string) ([]models.Withdrawal, error) {
	var withdrawals []models.Withdrawal

	rows, err := sq.Select(`"order"`, "sum", "processed_at").
		From("withdrawals").
		Where(sq.Eq{"login": login}).
		OrderBy("processed_at DESC").
		RunWith(s.DB).
		PlaceholderFormat(sq.Dollar).
		QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var withdrawal models.Withdrawal
		err = rows.Scan(&withdrawal.Order, &withdrawal.Sum, &withdrawal.ProcessedAt)
		if err != nil {
			return nil, err
		}
		withdrawals = append(withdrawals, withdrawal)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return withdrawals, nil
}
