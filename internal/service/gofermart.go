package service

import (
	"context"
	"fmt"
	"database/sql"
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
	Auth(ctx context.Context, login string) error
	PostOrders(ctx context.Context, login string, orderNum int) error
	GetOrders(ctx context.Context, login string) ([]models.Order, error)
	GetBalance(ctx context.Context, login string) (models.Balance, error)
	Withdraw(ctx context.Context, withdrawal models.Withdrawal) (models.Balance, error)
	Withdrawals(ctx context.Context, login string) ([]models.Withdrawal, error)
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
	ticker := time.NewTicker(1 * time.Second)
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
				order.ID = orderOld.ID
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
	var o models.Order
	url := fmt.Sprintf("%s/api/orders/%s", s.cfg.AccrualAddr, orderOld.Order)
	_, err := s.Client.R().SetContext(ctx).SetResult(&o).Get(url)
	if err != nil {
		return o, err
	}
	if o.Status == orderOld.Status {
		return o, err
	}
	o.Order = orderOld.Order
	return o, nil
}

func (s *Gofermart) updateStatus(ctx context.Context, orders []models.Order) error {
	tx, err := s.Storage.DB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = s.Storage.UpdateOrders(ctx, tx, orders)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
        return err
    }
	return nil
}
