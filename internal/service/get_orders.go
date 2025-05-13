package service

import (
	"context"
	"gofermart/internal/models"
)

func (s *Gofermart) GetOrders(ctx context.Context, userID string) ([]models.Order, error) {
	orders, err := s.Storage.GetOrders(ctx, userID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}