package storage

import (
	"context"
	"gofermart/internal/models"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func (s *Storage) GetBalance(ctx context.Context, userID string) (models.Balance, error) {
    var balance models.Balance
    var currentNull, withdrawnNull sql.NullFloat64
    
    query := `SELECT login, current, withdrawn FROM balances WHERE login = $1`
    row := s.DB.QueryRowContext(ctx, query, userID)
    
    err := row.Scan(&balance.ID, &currentNull, &withdrawnNull)
    if err != nil {
        if err.Error() == "sql: no rows in result set" {
            balance.ID = userID
            balance.Current = 0
            balance.Withdrawn = 0
            return balance, nil
        }
        return balance, err
    }

    if currentNull.Valid {
        balance.Current = currentNull.Float64
    }
    if withdrawnNull.Valid {
        balance.Withdrawn = withdrawnNull.Float64
    }

    return balance, nil
}
