package service

import (
	"context"
	"database/sql"
	"gofermart/internal/models"
)

func (s *Gofermart) Register(ctx context.Context, user models.User) error {
	tx, err := s.Storage.DB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if user.Login == "" || user.Password == "" {
		return ErrMalformedRequest
	}

	err = s.Storage.Register(ctx, tx, user)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
