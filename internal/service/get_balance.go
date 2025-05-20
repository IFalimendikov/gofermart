package service

import (
	"context"
	"gofermart/internal/models"
)

func (s *Gofermart) GetBalance(ctx context.Context, login string) (models.Balance, error) {
	balance, err := s.Storage.GetBalance(ctx, login)
	if err != nil {
		return balance, err
	}
	return balance, nil
}
