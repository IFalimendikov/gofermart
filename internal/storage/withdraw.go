package storage

import (
	"context"
	"database/sql"
	"gofermart/internal/models"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func (s *Storage) Withdraw(ctx context.Context, withdrawal models.Withdrawal) (models.Balance, error) {
	var balance sql.NullFloat64
	var accrual models.Balance

	tx, err := s.DB.Begin()
	if err != nil {
		return accrual, err
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(ctx, `SELECT current FROM balances WHERE login = $1`, withdrawal.ID)
	err = row.Scan(&balance)
	if err != nil {
		return accrual, err
	}

	if !balance.Valid || balance.Float64 < withdrawal.Sum {
		return accrual, ErrBalanceTooLow
	}

	_, err = tx.ExecContext(ctx,
		`UPDATE balances SET current = current - $1, withdrawn = withdrawn + $2 WHERE login = $3`,
		100, 100, withdrawal.ID)
	if err != nil {
		return accrual, err
	}

	_, err = tx.ExecContext(ctx,
		`INSERT into withdrawals ("order", login, sum, processed_at) VALUES ($1, $2, $3, $4)`,
		withdrawal.Order, withdrawal.ID, withdrawal.Sum, time.Now().Format(time.RFC3339))
	if err != nil {
		return accrual, err
	}

	var currentNull, withdrawnNull sql.NullFloat64
	row = tx.QueryRowContext(ctx, `SELECT current, withdrawn FROM balances WHERE login = $1`, withdrawal.ID)
	err = row.Scan(&currentNull, &withdrawnNull)
	if err != nil {
		return accrual, err
	}

	if currentNull.Valid {
		accrual.Current = currentNull.Float64
	}
	if withdrawnNull.Valid {
		accrual.Withdrawn = withdrawnNull.Float64
	}

	err = tx.Commit()
	if err != nil {
		return accrual, err
	}

	return accrual, nil
}
