package service

import (
	"context"
	"fmt"
	"gofermart/internal/config"
	"gofermart/internal/models"
	"gofermart/internal/storage"
	"log/slog"
	"time"

	"github.com/go-resty/resty/v2"
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

type Gofermart struct {
	Log     *slog.Logger
	cfg     *config.Config
	Storage *storage.Storage
	Client  *resty.Client
}

func New(log *slog.Logger, cfg config.Config, storage *storage.Storage, client *resty.Client) (*Gofermart, error) {
	service := Gofermart{
		Log:     log,
		cfg:     &cfg,
		Storage: storage,
		Client:  client,
	}
	return &service, nil
}

func (s *Gofermart) UpdateOrders(ctx context.Context) error {
	ticker := time.NewTicker(1 * time.Microsecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			var orders []models.Order
			ordersOld, err := s.getOrders(ctx)
			if err != nil {
				slog.Info("failed to get orders", "error", err)
				continue
			}

			for _, orderOld := range ordersOld {
				order, err := s.getStatus(ctx, orderOld)
				if err != nil {
					slog.Error("failed to get status", "error", err)
					continue
				}
				orders = append(orders, order)
			}
			
			err = s.updateStatus(ctx, orders)
			if err != nil {
				slog.Error("failed to update orders", "error", err)
				continue
			}
		}
	}
}

func (s *Gofermart) getOrders(ctx context.Context) ([]models.Order, error) {
	orders, err := s.Storage.GetOrdersNums(ctx)
	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, ErrNoNewAddresses
	}
	return orders, nil
}

func (s *Gofermart) getStatus(ctx context.Context, orderOld models.Order) (models.Order, error) {
	var order models.Order
	url := fmt.Sprintf("%s/api/orders/%s", s.cfg.AccrualAddr, orderOld.Order)
	_, err := s.Client.R().SetContext(ctx).SetResult(&order).Get(url)
	if err != nil {
		return order, err
	}
	if order.Status == orderOld.Status {
		return order, err
	}
	order.Order = orderOld.Order
	return order, nil
}

func (s *Gofermart) updateStatus(ctx context.Context, orders []models.Order) error {
	err := s.Storage.UpdateOrders(ctx, orders)
	if err != nil {
		return err
	}
	return nil
}
