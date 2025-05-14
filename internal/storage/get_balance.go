package storage

import (
	"context"
	"fmt"
	"gofermart/internal/models"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func (s *Storage) GetBalance(ctx context.Context, userID string) (models.Balance, error) {
	var balance models.Balance
	query := `SELECT login, current, withdrawn FROM balances WHERE login = $1`
	row  := s.DB.QueryRowContext(ctx, query, userID)
	
	err := row.Scan(&balance.ID, &balance.Current, &balance.Withdrawn) 
	if err != nil {
		return balance, err
	}

	fmt.Println("balance is")
	fmt.Println(balance)


	return balance, nil
}
