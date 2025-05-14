package storage

import (
	"context"
	"gofermart/internal/models"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func (s *Storage) Withdrawals(ctx context.Context, userID string) ([]models.Withdrawal, error) {
	var withdrawals []models.Withdrawal

	var query = `SELECT number, sum, processed_at FROM withdrawals WHERE login = $1 ORDER BY processed_at DESC`
	stmt, err := s.DB.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, userID)
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
	return withdrawals, err
}
