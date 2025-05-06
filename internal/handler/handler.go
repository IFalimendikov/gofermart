package handler

import (
	"context"
	"gofermart/internal/models"
	"log/slog"
)

type Service interface {
	Register(ctx context.Context, user models.User) error
	Login(ctx context.Context, user models.User) error
}

type Handler struct {
	Service Service
	log     *slog.Logger
}

func New(s Service, log *slog.Logger) *Handler {
	return &Handler{
		Service: s,
		log:     log,
	}
}
