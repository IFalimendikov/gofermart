package storage

import (
	"context"
	"time"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

func (s *Storage) PostOrders(ctx context.Context, login, orderNum string) error {
    var sUser string
    var sNumber string

    row := sq.Select("login" , "number").
        From("orders").
        Where(sq.Eq{"number": orderNum}).
        RunWith(s.DB).
        PlaceholderFormat(sq.Dollar).
        QueryRowContext(ctx)
    
    err := row.Scan(&sUser, &sNumber)
    if err != nil {
        if err == sql.ErrNoRows {
            _, err = sq.Insert("orders").
                Columns("number", "login", "status", "uploaded_at").
                Values(orderNum, login, "NEW", time.Now().UTC().Format(time.RFC3339)).
                RunWith(s.DB).
                PlaceholderFormat(sq.Dollar).
                ExecContext(ctx)
            return err
        }
        return err
    }

    switch {
    case login == sUser && orderNum == sNumber:
        return ErrDuplicateOrder
    case orderNum == sNumber && login != sUser:
        return ErrDuplicateNumber
    }

    return nil
}
