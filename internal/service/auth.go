package service

import (
	"context"
	"gofermart/internal/models"
)

func (s *Gofermart) Auth(ctx context.Context, userID string) error {
	user := models.User{
		ID: userID,
	}
	err := s.Storage.Auth(ctx, user)
	if err != nil {
		return err
	}
	return nil
}