package storage

import (
	"context"
	"gofermart/internal/models"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func (s *Storage) Withdraw(ctx context.Context, withdrawal models.Withdrawal) (models.Balance, error) {
	var balance int
	var accrual models.Balance

	tx, err := s.DB.Begin()
	if err != nil {
		return accrual, err
	}
	defer tx.Rollback()

	var queryBal = `SELECT current FROM balances WHERE user_id = $1`
	stmtBal, err := tx.PrepareContext(ctx, queryBal)
	if err != nil {
		return accrual, err
	}
	defer stmtBal.Close()

	row := tx.QueryRowContext(ctx, queryBal, withdrawal.ID)

	err = row.Scan(&balance)
	if err != nil {
		return accrual, err
	}

	if balance < withdrawal.Sum {
		return accrual, ErrBalanceTooLow
	}

	var queryNewBal = `UPDATE balances SET current = current - $1, withdrawn = withdrawn + $2 WHERE user_id = $3`
	stmtNewBal, err := tx.PrepareContext(ctx, queryNewBal)
	if err != nil {
		return accrual, err
	}
	defer stmtNewBal.Close()

	_, err = tx.ExecContext(ctx, queryNewBal, withdrawal.Sum, withdrawal.Sum, withdrawal.ID)
	if err != nil {
		return accrual, err
	}

	var queryWihdraw = `INSERT into withdrawals (order, user_id, sum, processed_at) VALUES ($1, $2, $3, $4)`
	stmtWihdraw, err := tx.PrepareContext(ctx, queryWihdraw)
	if err != nil {
		return accrual, err
	}
	defer stmtWihdraw.Close()

	_, err = tx.ExecContext(ctx, queryWihdraw, withdrawal.Order, withdrawal.ID, withdrawal.Sum, time.Now().Format(time.RFC3339))
	if err != nil {
		return accrual, err
	}

	var queryAccrual = `SELECT current, withdrawn FROM balances WHERE user_id = $1`
	stmtAccrual, err := tx.PrepareContext(ctx, queryAccrual)
	if err != nil {
		return accrual, err
	}
	defer stmtAccrual.Close()

	row = tx.QueryRowContext(ctx, queryAccrual, withdrawal.ID)

	err = row.Scan(&accrual.Current, &accrual.Withdrawn)
	if err != nil {
		return accrual, err
	}

	err = tx.Commit()
	if err != nil {
		return accrual, err
	}
	return accrual, nil
}
