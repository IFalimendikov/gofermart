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
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var queryUser = `INSERT into users (login, password) VALUES ($1, $2)`
	stmtUser, err := tx.PrepareContext(ctx, queryUser)
	if err != nil {
		return err
	}
	defer stmtUser.Close()

	_, err = tx.ExecContext(ctx, queryUser, user.Login, user.Password)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return ErrDuplicateLogin
		}
		return err
	}

	var queryBal = `INSERT into balances (login) VALUES ($1)`
	stmtBal, err := tx.PrepareContext(ctx, queryBal)
	if err != nil {
		return err
	}
	defer stmtBal.Close()

	_, err = tx.ExecContext(ctx, queryBal, user.Login)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return err
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}