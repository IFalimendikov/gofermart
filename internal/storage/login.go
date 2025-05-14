package storage

import (
	"context"
	"database/sql"
	"errors"
	"gofermart/internal/models"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func (s *Storage) Login(ctx context.Context, user models.User) error {
    checkQuery := `SELECT login FROM users WHERE login = $1 AND password = $2`
    
    var login string
    err := s.DB.QueryRowContext(ctx, checkQuery, user.Login, user.Password).Scan(&login)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return ErrWrongPassword
        }
        return err
    }

    return nil
}