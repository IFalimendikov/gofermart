package storage

import (
	"context"
	"database/sql"
	"gofermart/internal/models"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func (s *Storage) GetBalance(ctx context.Context, login string) (models.Balance, error) {
	var balance models.Balance
	var currentNull, withdrawnNull sql.NullFloat64

	query := `SELECT login, current, withdrawn FROM balances WHERE login = $1`
	row := s.DB.QueryRowContext(ctx, query, login)

	err := row.Scan(&balance.ID, &currentNull, &withdrawnNull)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			balance.ID = login
			balance.Current = 0
			balance.Withdrawn = 0
			return balance, nil
		}
		return balance, err
	}

	if currentNull.Valid {
		balance.Current = currentNull.Float64
	} else {
		balance.Current = 0.0
	}

	if withdrawnNull.Valid {
		balance.Withdrawn = withdrawnNull.Float64
	} else {
		balance.Withdrawn = 0.0
	}

	return balance, nil
}
