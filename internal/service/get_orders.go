package service

import (
	"context"
	"gofermart/internal/models"
)

func (s *Gofermart) GetOrders(ctx context.Context, login string) ([]models.Order, error) {
	orders, err := s.Storage.GetOrders(ctx, login)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
