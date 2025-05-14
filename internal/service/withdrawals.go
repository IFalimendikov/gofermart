package service

import (
	"context"
	"gofermart/internal/models"
)

func (s *Gofermart) Withdrawals(ctx context.Context, userID string) ([]models.Withdrawal, error) {
	var withdrawals []models.Withdrawal

	withdrawals, err := s.Storage.Withdrawals(ctx, userID)
	if err != nil {
		return nil, err
	}
	return withdrawals, nil
}
