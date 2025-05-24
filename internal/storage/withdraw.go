package storage

import (
	"context"
	"database/sql"
	"gofermart/internal/models"
	"time"

	sq "github.com/Masterminds/squirrel"
)

func (s *Storage) Withdraw(ctx context.Context, withdrawal models.Withdrawal) (models.Balance, error) {
	var balance sql.NullFloat64
	var accrual models.Balance

	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
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
		withdrawal.Sum, withdrawal.Sum, withdrawal.ID)
	if err != nil {
		return accrual, err
	}

	_, err = sq.Insert("withdrawals").
		Columns(`"order"`, "login", "sum", "processed_at").
		Values(withdrawal.Order, withdrawal.ID, withdrawal.Sum, time.Now().Format(time.RFC3339)).
		RunWith(tx).
		PlaceholderFormat(sq.Dollar).
		ExecContext(ctx)
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
	} else {
		accrual.Current = 0.0
	}

	if withdrawnNull.Valid {
		accrual.Withdrawn = withdrawnNull.Float64
	} else {
		accrual.Withdrawn = 0.0
	}

	err = tx.Commit()
	if err != nil {
		return accrual, err
	}

	return accrual, nil
}
