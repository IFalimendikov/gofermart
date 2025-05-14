package storage

import (
	"context"
	"log"
	"gofermart/internal/models"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func (s *Storage) Withdraw(ctx context.Context, withdrawal models.Withdrawal) (models.Balance, error) {
	log.Printf("Starting withdrawal process for user ID: %s, amount: %f, order: %s", withdrawal.ID, withdrawal.Sum, withdrawal.Order)
	
	var balance float64
	var accrual models.Balance

	tx, err := s.DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return accrual, err
	}
	defer tx.Rollback()

	var queryBal = `SELECT current FROM balances WHERE login = $1`
	stmtBal, err := tx.PrepareContext(ctx, queryBal)
	if err != nil {
		log.Printf("Error preparing balance query: %v", err)
		return accrual, err
	}
	defer stmtBal.Close()

	row := tx.QueryRowContext(ctx, queryBal, withdrawal.ID)

	err = row.Scan(&balance)
	if err != nil {
		log.Printf("Error scanning balance for user %s: %v", withdrawal.ID, err)
		return accrual, err
	}
	log.Printf("Current balance for user %s: %f", withdrawal.ID, balance)

	if balance < withdrawal.Sum {
		log.Printf("Insufficient balance for user %s. Required: %f, Available: %f", withdrawal.ID, withdrawal.Sum, balance)
		return accrual, ErrBalanceTooLow
	}

	var queryNewBal = `UPDATE balances SET current = current - $1, withdrawn = withdrawn + $2 WHERE login = $3`
	stmtNewBal, err := tx.PrepareContext(ctx, queryNewBal)
	if err != nil {
		log.Printf("Error preparing update balance query: %v", err)
		return accrual, err
	}
	defer stmtNewBal.Close()

	_, err = tx.ExecContext(ctx, queryNewBal, withdrawal.Sum, withdrawal.Sum, withdrawal.ID)
	if err != nil {
		log.Printf("Error updating balance: %v", err)
		return accrual, err
	}
	log.Printf("Successfully updated balance for user %s", withdrawal.ID)

	var queryWihdraw = `INSERT into withdrawals (number, login, sum, processed_at) VALUES ($1, $2, $3, $4)`
	stmtWihdraw, err := tx.PrepareContext(ctx, queryWihdraw)
	if err != nil {
		log.Printf("Error preparing withdrawal insert query: %v", err)
		return accrual, err
	}
	defer stmtWihdraw.Close()

	_, err = tx.ExecContext(ctx, queryWihdraw, withdrawal.Order, withdrawal.ID, withdrawal.Sum, time.Now().Format(time.RFC3339))
	if err != nil {
		log.Printf("Error inserting withdrawal record: %v", err)
		return accrual, err
	}
	log.Printf("Successfully recorded withdrawal for user %s", withdrawal.ID)

	var queryAccrual = `SELECT current, withdrawn FROM balances WHERE login = $1`
	stmtAccrual, err := tx.PrepareContext(ctx, queryAccrual)
	if err != nil {
		log.Printf("Error preparing final balance query: %v", err)
		return accrual, err
	}
	defer stmtAccrual.Close()

	row = tx.QueryRowContext(ctx, queryAccrual, withdrawal.ID)

	err = row.Scan(&accrual.Current, &accrual.Withdrawn)
	if err != nil {
		log.Printf("Error scanning final balance: %v", err)
		return accrual, err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		return accrual, err
	}
	
	log.Printf("Withdrawal completed successfully for user %s. New balance: %f, Total withdrawn: %f", 
		withdrawal.ID, accrual.Current, accrual.Withdrawn)
	return accrual, nil
}