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
    var existingUser models.User
    err := s.DB.QueryRowContext(ctx, 
        "SELECT user_id, login, password FROM users WHERE login = $1", 
        user.Login).Scan(&existingUser.ID, &existingUser.Login, &existingUser.Password)
    
    if err == nil {
        if existingUser.ID == user.ID && existingUser.Password == user.Password {
            return nil
        }
        return ErrDuplicateLogin
    }

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var queryUser = `INSERT into users (user_id, login, password) VALUES ($1, $2, $3)`
	stmtUser, err := tx.PrepareContext(ctx, queryUser)
	if err != nil {
		return err
	}
	defer stmtUser.Close()

	_, err = tx.ExecContext(ctx, queryUser, user.ID, user.Login, user.Password)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return err
		}
		return err
	}

	var queryBal = `INSERT into balances (user_id) VALUES ($1)`
	stmtBal, err := tx.PrepareContext(ctx, queryBal)
	if err != nil {
		return err
	}
	defer stmtBal.Close()

	_, err = tx.ExecContext(ctx, queryBal, user.ID)
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
