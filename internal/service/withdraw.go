package service

import (
	"context"
	"database/sql"
	"github.com/ShiraazMoollatjie/goluhn"
	"gofermart/internal/models"
)

func (s *Gofermart) Withdraw(ctx context.Context, withdrawal models.Withdrawal) (models.Balance, error) {
	var balance models.Balance
	tx, err := s.Storage.DB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return balance, err
	}
	defer tx.Rollback()

	if withdrawal.ID == "" || withdrawal.Order == "" {
		return balance, ErrWrongFormat
	}

	err = goluhn.Validate(withdrawal.Order)
	if err != nil {
		return balance, ErrWrongFormat
	}

	balance, err = s.Storage.Withdraw(ctx, tx, withdrawal)
	if err != nil {
		return balance, err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return models.Balance{}, err
	}
	return balance, nil
}
