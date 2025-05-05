package service

import (
	"context"
	"gofermart/internal/models"

	"github.com/deatil/go-encoding/base62"
)

func (s *Service) Register(ctx context.Context, user models.User) error {
	err := s.Storage.Register(ctx, user)
	if err != nil {
		return err
	}
}