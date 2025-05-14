package service

import (
	"context"
	"gofermart/internal/models"
)

func (s *Gofermart) GetBalance(ctx context.Context, userID string) (models.Balance, error) {
	balance, err := s.Storage.GetBalance(ctx, userID)
	if err != nil {
		return balance, err
	}
	return balance, nil
}