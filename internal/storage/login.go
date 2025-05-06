package storage

import (
	"context"
	"gofermart/internal/models"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func (s *Storage) Login(ctx context.Context, user models.User) error {
	var query = `UPDATE users SET connected = true WHERE user_id = $1 AND login = $2 AND password = $3`
	result, err := s.DB.ExecContext(ctx, query, user.ID, user.Login, user.Password)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrWrongPassword
	}
	return nil
}
