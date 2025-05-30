package storage

import (
	"context"
	"gofermart/internal/models"
	"time"

	sq "github.com/Masterminds/squirrel"
)

func (s *Storage) Withdraw(ctx context.Context, runner sq.BaseRunner, withdrawal models.Withdrawal) (models.Balance, error) {
	var balance float64
	var accrual models.Balance

	row := sq.Select("current").
		From("balances").
		Where(sq.Eq{"login": withdrawal.ID}).
		RunWith(runner).
		PlaceholderFormat(sq.Dollar).
		QueryRowContext(ctx)

	err := row.Scan(&balance)
	if err != nil {
		return accrual, err
	}

	if balance < withdrawal.Sum {
		return accrual, ErrBalanceTooLow
	}

	_, err = sq.Update("balances").
		Set("current", sq.Expr("current - $1", withdrawal.Sum)).
		Set("withdrawn", sq.Expr("withdrawn + $1", withdrawal.Sum)).
		Where(sq.Eq{"login": withdrawal.ID}).
		RunWith(runner).
		PlaceholderFormat(sq.Dollar).
		ExecContext(ctx)
	if err != nil {
		return accrual, err
	}

	_, err = sq.Insert("withdrawals").
		Columns(`"order"`, "login", "sum", "processed_at").
		Values(withdrawal.Order, withdrawal.ID, withdrawal.Sum, time.Now().UTC().Format(time.RFC3339)).
		RunWith(runner).
		PlaceholderFormat(sq.Dollar).
		ExecContext(ctx)
	if err != nil {
		return accrual, err
	}

	row = sq.Select("current" , "withdrawn").
		From("balances").
		Where(sq.Eq{"login": withdrawal.ID}).
		RunWith(runner).
		PlaceholderFormat(sq.Dollar).
		QueryRowContext(ctx)

	err = row.Scan(&accrual.Current, &accrual.Withdrawn)
	if err != nil {
		return accrual, err
	}

	return accrual, nil
}
