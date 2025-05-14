package service

import (
	"context"
	"gofermart/internal/models"
	"github.com/ShiraazMoollatjie/goluhn"
)

func (s *Gofermart) Withdraw(ctx context.Context, withdrawal models.Withdrawal) (models.Balance, error) {
	var balance models.Balance
	err := goluhn.Validate(withdrawal.Order)
	if err != nil {
		return balance, ErrWrongFormat
	}

	balance, err = s.Storage.Withdraw(ctx, withdrawal)
	if err != nil {
		return balance, err
	}
	return balance, nil
}
