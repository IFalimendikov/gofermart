package service

import (
	"context"
	"log/slog"
	"gofermart/internal/models"
	"gofermart/internal/storage"
)

type Service interface {
	Register(ctx context.Context, user models.User) error
	Login(ctx context.Context, user models.User) error
}

type Gofermart struct {
	Log *slog.Logger
	Storage *storage.Storage
}

func New(log *slog.Logger, storage *storage.Storage) (*Gofermart, error){
	service := Gofermart {
		Log: log,
		Storage: storage,
	}
	return &service, nil
}