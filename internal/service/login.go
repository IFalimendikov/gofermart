package service

import (
	"context"
	"gofermart/internal/models"
)

func (s *Gofermart) Login(ctx context.Context, user models.User) error {
	err := s.Storage.Login(ctx, user)
	if err != nil {
		return err
	}
	return nil
}