package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gofermart/internal/config"
	"gofermart/internal/models"
	"gofermart/internal/service"
	"gofermart/internal/storage"
	"log/slog"
)

func setupWithdrawalsOrdersTestDB(t *testing.T) *sql.DB {
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
	_, err = db.Exec(storage.WithdrawalsQuery)
	require.NoError(t, err)
	return db
}

func TestGofermart_Withdrawals(t *testing.T) {
	db := setupWithdrawalsOrdersTestDB(t)
	defer db.Close()
	defer func() {
		_, err := db.Exec(`DROP TABLE IF EXISTS withdrawals`)
		require.NoError(t, err)
		_, err = db.Exec(`DROP TABLE IF EXISTS balances`)
		require.NoError(t, err)
		_, err = db.Exec(`DROP TABLE IF EXISTS users`)
		require.NoError(t, err)
		_, err = db.Exec(`DROP TABLE IF EXISTS orders`)
		require.NoError(t, err)
	}()

	storage := &storage.Storage{DB: db}
	service := &service.Gofermart{
		Storage: storage,
		Log:     slog.Default(),
	}
	handler := &Handler{
		Service: service,
	}

	testUser := models.User{
		Login:    "testuser",
		Password: "testpass",
	}
	err := service.Register(context.Background(), testUser)
	require.NoError(t, err)

	timeFormat := "2006-01-02 15:04:05.999999999-07:00"
	now := time.Now().UTC().Format(timeFormat)
	later := time.Now().UTC().Add(time.Hour).Format(timeFormat)

	testWithdrawals := []struct {
		order       string
		sum         float64
		processedAt string
	}{
		{
			order:       "79927398713",
			sum:         100.0,
			processedAt: now,
		},
		{
			order:       "79927398714",
			sum:         200.0,
			processedAt: later,
		},
	}

	for _, w := range testWithdrawals {
		_, err = db.Exec(`
			INSERT INTO withdrawals (login, "order", sum, processed_at)
			VALUES ($1, $2, $3, $4)
		`, testUser.Login, w.order, w.sum, w.processedAt)
		require.NoError(t, err)
	}

	tests := []struct {
		name       string
		login      string
		wantStatus int
		wantLen    int
	}{
		{
			name:       "valid user with withdrawals",
			login:      "testuser",
			wantStatus: http.StatusOK,
			wantLen:    2,
		},
		{
			name:       "user with no withdrawals",
			login:      "nonexistentuser",
			wantStatus: http.StatusNoContent,
			wantLen:    0,
		},
	}

	gin.SetMode(gin.TestMode)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(func(c *gin.Context) {
				c.Set("login", tt.login)
				c.Next()
			})

			router.GET("/api/user/withdrawals", func(c *gin.Context) {
				handler.Withdrawals(c)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusOK {
				var gotWithdrawals []models.Withdrawal
				err = json.Unmarshal(w.Body.Bytes(), &gotWithdrawals)
				require.NoError(t, err)
				assert.Equal(t, tt.wantLen, len(gotWithdrawals))

				if len(gotWithdrawals) > 1 {
					for i := 0; i < len(gotWithdrawals)-1; i++ {

						t1, err := time.Parse(timeFormat, gotWithdrawals[i].ProcessedAt)
						require.NoError(t, err)
						t2, err := time.Parse(timeFormat, gotWithdrawals[i+1].ProcessedAt)
						require.NoError(t, err)
						
						assert.True(t, t1.After(t2), "Withdrawals should be ordered by processed_at DESC")
					}
				}
			}
		})
	}
}
