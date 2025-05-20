package handler

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"testing"

	"gofermart/internal/config"
	"gofermart/internal/models"
	"gofermart/internal/service"
	"gofermart/internal/storage"
	"log/slog"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupGetOrdersTestDB(t *testing.T) *sql.DB {
	cfg := config.Config{}
	
	err := config.Read(&cfg)
	require.NoError(t, err, "Failed to read config")
	db, err := sql.Open("postgres", cfg.DatabaseURI)
	require.NoError(t, err)
	err = db.Ping()
	require.NoError(t, err)

	_, err = db.Exec(storage.OrdersQuery)
	require.NoError(t, err)
	_, err = db.Exec(storage.UsersQuery)
	require.NoError(t, err)
	_, err = db.Exec(storage.BalancesQuery)
	require.NoError(t, err)
	return db
}

func TestGofermart_ConcurrentBalanceOperations(t *testing.T) {
	db := setupGetOrdersTestDB(t)
	defer db.Close()
	defer func() {
		_, err := db.Exec(`DROP TABLE IF EXISTS orders`)
		require.NoError(t, err)
		_, err = db.Exec(`DROP TABLE IF EXISTS users`)
		require.NoError(t, err)
		_, err = db.Exec(`DROP TABLE IF EXISTS balances`)
		require.NoError(t, err)
	}()

	storage := &storage.Storage{DB: db}
	service := &service.Gofermart{
		Storage: storage,
		Log:     slog.Default(),
	}

	testUser := models.User{
		Login:    "testuser",
		Password: "testpass",
	}
	err := service.Register(context.Background(), testUser)
	require.NoError(t, err, "Failed to register test user")

	var exists bool
	err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM balances WHERE login = $1)`, testUser.Login).Scan(&exists)
	require.NoError(t, err)

	if exists {
		_, err = db.Exec(`
			UPDATE balances 
			SET current = $1, withdrawn = $2
			WHERE login = $3
		`, 1000.0, 0.0, testUser.Login)
	} else {
		_, err = db.Exec(`
			INSERT INTO balances (login, current, withdrawn)
			VALUES ($1, $2, $3)
		`, testUser.Login, 1000.0, 0.0)
	}
	require.NoError(t, err)

	var wg sync.WaitGroup
	concurrentOperations := 10
	wg.Add(concurrentOperations)

	for i := 0; i < concurrentOperations; i++ {
		go func(i int) {
			defer wg.Done()

			tx, err := db.Begin()
			require.NoError(t, err)

			defer func() {
				if err != nil {
					tx.Rollback()
				}
			}()

			if i%2 == 0 {
				orderNum := fmt.Sprintf("7992739871%d", i)
				_, err := tx.Exec(`
					INSERT INTO orders (number, login, status, accrual, uploaded_at)
					VALUES ($1, $2, $3, $4, NOW())
					ON CONFLICT (number) DO NOTHING
				`, orderNum, testUser.Login, "PROCESSED", 50.0)
				require.NoError(t, err)

				_, err = tx.Exec(`
					UPDATE balances 
					SET current = current + 50.0
					WHERE login = $1
				`, testUser.Login)
				require.NoError(t, err)
			} else {
				_, err := tx.Exec(`
					UPDATE balances 
					SET current = current - 25.0,
					    withdrawn = withdrawn + 25.0
					WHERE login = $1
				`, testUser.Login)
				require.NoError(t, err)
			}

			err = tx.Commit()
			require.NoError(t, err)
		}(i)
	}

	wg.Wait()

	var currentBalance, withdrawnBalance float64
	err = db.QueryRow(`
		SELECT current, withdrawn FROM balances WHERE login = $1
	`, testUser.Login).Scan(&currentBalance, &withdrawnBalance)
	require.NoError(t, err)

	assert.Equal(t, 1125.0, currentBalance, "Unexpected final current balance")
	assert.Equal(t, 125.0, withdrawnBalance, "Unexpected final withdrawn balance")
}
