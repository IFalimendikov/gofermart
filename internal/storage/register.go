package storage

import (
	"context"
	"errors"
	"gofermart/internal/models"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func (s *Storage) Register(ctx context.Context, user models.User) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

    _, err = sq.Insert("users").
        Columns("login", "password").
        Values(user.Login, user.Password).
        RunWith(tx).
        PlaceholderFormat(sq.Dollar).
        ExecContext(ctx)
    if err != nil {
        var pgErr *pgconn.PgError
        if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
            return ErrDuplicateLogin
        }
        return err
    }

	_, err = sq.Insert("balances").
		Columns("login").
		Values(user.Login).
		RunWith(tx).
		PlaceholderFormat(sq.Dollar).
		ExecContext(ctx)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return err
		}
		return err
	}

	return tx.Commit()
}
