package handler

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

    "github.com/brianvoe/gofakeit/v7"
    "github.com/ShiraazMoollatjie/goluhn"
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

func setupWithdrawOrdersTestDB(t *testing.T) *sql.DB {
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

func TestGofermart_Withdraw(t *testing.T) {
    gofakeit.Seed(0)

    db := setupWithdrawOrdersTestDB(t)
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
        Login:    gofakeit.Username(),
        Password: gofakeit.Password(true, true, true, true, false, 10),
    }
    
    initialBalance := gofakeit.Float64Range(1000, 5000)

    err := service.Register(context.Background(), testUser)
    require.NoError(t, err)
    
    _, err = db.Exec(`
        INSERT INTO balances (login, current, withdrawn)
        VALUES ($1, $2, $3)
        ON CONFLICT (login) DO UPDATE 
        SET current = $2, withdrawn = $3
    `, testUser.Login, initialBalance, 0.0)
    require.NoError(t, err)

    validLuhnNumber := goluhn.Generate(8)
    
    tests := []struct {
        name       string
        withdrawal models.Withdrawal
        login      string
        wantStatus int
        wantBalance models.Balance
    }{
        {
            name: "valid withdrawal",
            withdrawal: models.Withdrawal{
                Order: validLuhnNumber,
                Sum:   gofakeit.Float64Range(100, initialBalance/2),
            },
            login:      testUser.Login,
            wantStatus: http.StatusOK,

        },
        {
            name: "invalid luhn number",
            withdrawal: models.Withdrawal{
                Order: gofakeit.Numerify("#####"),
                Sum:   gofakeit.Float64Range(50, 100),
            },
            login:      testUser.Login,
            wantStatus: http.StatusUnprocessableEntity,
        },
        {
            name: "insufficient balance",
            withdrawal: models.Withdrawal{
                Order: validLuhnNumber,
                Sum:   initialBalance * gofakeit.Float64Range(1.1, 2.0),
            },
            login:      testUser.Login,
            wantStatus: http.StatusBadRequest,
        },
        {
            name: "random valid withdrawal",
            withdrawal: models.Withdrawal{
                Order: validLuhnNumber,
                Sum:   gofakeit.Float64Range(1, initialBalance),
            },
            login:      testUser.Login,
            wantStatus: http.StatusOK,
        },
    }
    
    gin.SetMode(gin.TestMode)
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if tt.name != "valid withdrawal" {
                _, err = db.Exec(`
                    UPDATE balances
                    SET current = $1, withdrawn = $2
                    WHERE login = $3
                `, 1000.0, 0.0, testUser.Login)
                require.NoError(t, err)
            }
            
            router := gin.New()
            router.Use(func(c *gin.Context) {
                c.Set("login", tt.login)
                c.Next()
            })
            
            router.POST("/api/user/balance/withdraw", func(c *gin.Context) {
                handler.Withdraw(c)
            })
            
            withdrawalJSON, err := json.Marshal(tt.withdrawal)
            require.NoError(t, err)
            
            w := httptest.NewRecorder()
            req := httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", bytes.NewBuffer(withdrawalJSON))
            req.Header.Set("Content-Type", "application/json")
            
            router.ServeHTTP(w, req)
            
            assert.Equal(t, tt.wantStatus, w.Code)
            
            if tt.wantStatus == http.StatusOK {
                var gotBalance models.Balance
                err = json.Unmarshal(w.Body.Bytes(), &gotBalance)
                require.NoError(t, err)
                assert.Equal(t, tt.wantBalance.Current, gotBalance.Current)
                assert.Equal(t, tt.wantBalance.Withdrawn, gotBalance.Withdrawn)
            }
        })
    }
}