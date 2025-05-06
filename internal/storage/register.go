package storage

import (
	"context"
	"errors"
	"gofermart/internal/models"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func (s *Storage) Register(ctx context.Context, user models.User) error {
	var query = `INSERT into users (user_id, login, password) VALUES ($1, $2, $3)`
	_, err := s.DB.ExecContext(ctx, query, user.ID, user.Login, user.Password)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err,&pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return ErrDuplicateLogin
		}
		return err
	}
	return nil
}
