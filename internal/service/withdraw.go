package service

import (
	"context"
	"github.com/ShiraazMoollatjie/goluhn"
	"gofermart/internal/models"
)

func (s *Gofermart) Withdraw(ctx context.Context, withdrawal models.Withdrawal) (models.Balance, error) {
	var balance models.Balance
	if withdrawal.ID == "" || withdrawal.Order == "" {
		return balance, ErrWrongFormat
	}

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
