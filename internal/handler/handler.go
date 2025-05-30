package handler

import (
	"context"
	"gofermart/internal/config"
	"gofermart/internal/models"
	"log/slog"
)

type Service interface {
	Register(ctx context.Context, user models.User) error
	Login(ctx context.Context, user models.User) error
	PostOrders(ctx context.Context, login string, orderNum int) error
	GetOrders(ctx context.Context, login string) ([]models.Order, error)
	GetBalance(ctx context.Context, login string) (models.Balance, error)
	Withdraw(ctx context.Context, withdrawal models.Withdrawal) (models.Balance, error)
	Withdrawals(ctx context.Context, login string) ([]models.Withdrawal, error)
}

type Handler struct {
	Service Service
	log     *slog.Logger
	cfg     *config.Config
}

func New(s Service, log *slog.Logger, cfg *config.Config) *Handler {
	return &Handler{
		Service: s,
		log:     log,
		cfg:     cfg,
	}
}
