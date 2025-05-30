package storage

import (
	"context"
	"gofermart/internal/models"

	sq "github.com/Masterminds/squirrel"
)

func (s *Storage) GetBalance(ctx context.Context, login string) (models.Balance, error) {
	var balance models.Balance

	row, err := sq.Select("login", "current", "withdrawn").
		From("balances").
		Where(sq.Eq{"login": login}).
		RunWith(s.DB).
		PlaceholderFormat(sq.Dollar).
		QueryContext(ctx)
	if err != nil {
		return balance, err
	}

	err = row.Scan(&balance.ID, &balance.Current, &balance.Withdrawn)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			balance.ID = login
			balance.Current = 0
			balance.Withdrawn = 0
			return balance, nil
		}
		return balance, err
	}
	if err = row.Err(); err != nil {
		return balance, err
	}

	return balance, nil
}
