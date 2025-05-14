package storage

import (
	"context"
	"database/sql"
	"fmt"

	"gofermart/internal/models"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func (s *Storage) GetOrders(ctx context.Context, userID string) ([]models.Order, error) {
    orders := make([]models.Order, 0)
    query := `SELECT number, status, accrual, uploaded_at FROM orders WHERE login = $1 ORDER BY uploaded_at DESC`
    rows, err := s.DB.Query(query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var order models.Order
        var accrual sql.NullFloat64
        err := rows.Scan(&order.Order, &order.Status, &accrual, &order.UploadedAt)
        if err != nil {
            return nil, err
        }
        if accrual.Valid {
            order.Accrual = accrual.Float64
        }
        orders = append(orders, order)
    }
    if len(orders) == 0 {
        return nil, ErrNoOrdersFound
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    // Print orders to console
    fmt.Printf("Found %d orders for user %s:\n", len(orders), userID)
    for i, order := range orders {
        fmt.Printf("Order %d:\n", i+1)
        fmt.Printf("  Order Number: %s\n", order.Order)
        fmt.Printf("  Status: %s\n", order.Status)
        fmt.Printf("  Accrual: %.2f\n", order.Accrual)
        fmt.Printf("  Uploaded At: %s\n", order.UploadedAt)
        fmt.Println()
    }

    return orders, nil
}
