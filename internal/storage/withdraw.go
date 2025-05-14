package storage

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"gofermart/internal/models"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func (s *Storage) Withdraw(ctx context.Context, withdrawal models.Withdrawal) (models.Balance, error) {
	log.Printf("Starting withdrawal process for user %s, amount: %f, order: %s", withdrawal.ID, withdrawal.Sum, withdrawal.Order)
	
	var balance sql.NullFloat64
	var accrual models.Balance

	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		return accrual, err
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(ctx, `SELECT current FROM balances WHERE login = $1`, withdrawal.ID)
	err = row.Scan(&balance)
	if err != nil {
		log.Printf("Failed to fetch current balance for user %s: %v", withdrawal.ID, err)
		return accrual, err
	}
	log.Printf("Current balance for user %s: %f", withdrawal.ID, balance.Float64)

	if !balance.Valid || balance.Float64 < withdrawal.Sum {
		log.Printf("Insufficient balance for user %s. Required: %f, Available: %f", 
			withdrawal.ID, withdrawal.Sum, balance.Float64)
		return accrual, ErrBalanceTooLow
	}

	_, err = tx.ExecContext(ctx,
		`UPDATE balances SET current = current - $1, withdrawn $2 WHERE login = $3`,
		withdrawal.Sum, withdrawal.Sum, withdrawal.ID)
	if err != nil {
		log.Printf("Failed to update balance for user %s: %v", withdrawal.ID, err)
		return accrual, err
	}
	log.Printf("Successfully updated balance for user %s", withdrawal.ID)

	_, err = tx.ExecContext(ctx,
		`INSERT into withdrawals ("order", login, sum, processed_at) VALUES ($1, $2, $3, $4)`,
		withdrawal.Order, withdrawal.ID, withdrawal.Sum, time.Now().Format(time.RFC3339))
	if err != nil {
		log.Printf("Failed to insert withdrawal record: %v", err)
		return accrual, err
	}
	log.Printf("Successfully recorded withdrawal transaction")

	var currentNull, withdrawnNull sql.NullFloat64
	row = tx.QueryRowContext(ctx, `SELECT current, withdrawn FROM balances WHERE login = $1`, withdrawal.ID)
	err = row.Scan(&currentNull, &withdrawnNull)
	if err != nil {
		log.Printf("Failed to fetch updated balance: %v", err)
		return accrual, err
	}

	if currentNull.Valid {
		accrual.Current = currentNull.Float64
	}
	if withdrawnNull.Valid {
		accrual.Withdrawn = withdrawnNull.Float64
	}
	log.Printf("Final balance - Current: %f, Withdrawn: %f", accrual.Current, accrual.Withdrawn)

	err = tx.Commit()
	if err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return accrual, err
	}
	log.Printf("Successfully completed withdrawal process")

	return accrual, errors.New("Biba")
}
