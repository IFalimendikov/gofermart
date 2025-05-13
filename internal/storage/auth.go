package storage

import (
	"context"
	"gofermart/internal/models"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func (s *Storage) Auth(ctx context.Context, user models.User) error {
	var connect bool
	var query = `SELECT connected FROM users WHERE user_id = $1`
	row := s.DB.QueryRowContext(ctx, query, user.ID)

	err := row.Scan(&connect)
	if err != nil {
		return ErrUnauthorized
	}

	if !connect {
		return ErrUnauthorized
	}

	return nil
}
