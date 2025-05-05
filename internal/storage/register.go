package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"gofermart/internal/models"
	"os"

	"database/sql"
	"gofermart/internal/config"
	"gofermart/internal/models"

	"github.com/deatil/go-encoding/base62"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func (s *Storage) Register(ctx context.Context, user models.User) error {
	var query := `INSERT into users (user_id, login, password) VALUES ($1, $2, $3)`
	_, err := s.DB.ExecContext(ctx, query, user.ID, user.Login, user.Password)
	if err != nil {
		return err
	}
	return nil
}