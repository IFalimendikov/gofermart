package service

import (
	"context"
	"gofermart/internal/models"
)

func (s *Gofermart) Register(ctx context.Context, user models.User) error {
	err := s.Storage.Register(ctx, user)
	if err != nil {
		return err
	}
	return nil
}