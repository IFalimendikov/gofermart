package storage

import (
	"context"

	"gofermart/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func (s *Storage) GetBalance(ctx context.Context, userID string) (models.Balance, error) {
	var balance models.Balance
	query := `SELECT user_id, current, withdrawn FROM balances WHERE user_id = $1`
	row  := s.DB.QueryRowContext(ctx, query, userID)
	
	err := row.Scan(&balance) 
	if err != nil {
		return balance, err
	}

	return balance, nil
}
