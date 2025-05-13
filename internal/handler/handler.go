package handler

import (
	"context"
	"gofermart/internal/models"
	"log/slog"
)

type Service interface {
	Register(ctx context.Context, user models.User) error
	Login(ctx context.Context, user models.User) error
	Auth(ctx context.Context, userID string) error
	PostOrders(ctx context.Context, userID string, orderNum int) error
	GetOrders(ctx context.Context, userID string) ([]models.Order, error)
	GetBalance(ctx context.Context, userID string) (models.Balance, error)
	Withdraw(ctx context.Context, withdrawal models.Withdrawal) (models.Balance, error)
	Withdrawals(ctx context.Context, userID string) ([]models.Withdrawal, error)
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
