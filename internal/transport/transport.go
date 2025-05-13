package transport

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"gofermart/internal/config"
	"gofermart/internal/handler"
	"gofermart/internal/models"
	"gofermart/internal/storage"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Service interface {
	Register(ctx context.Context, user models.User) error
	Login(ctx context.Context, user models.User) error
	Auth(ctx context.Context, userID string) error
	PostOrders(ctx context.Context, userID string, orderNum int) error
	GetOrders(ctx context.Context, userID string) ([]models.Order, error)
	GetBalance(ctx context.Context, userID string) (models.Balance, error)
	Withdraw(ctx context.Context, withdrawal models.Withdrawal) (models.Balance, error)
	Withdrawals(ctx context.Context, userID string) ([]models.Withdrawal, error)
}

type Transport struct {
	Handler *handler.Handler
	Log     *slog.Logger
	Config  *config.Config
}

type Claim struct {
	jwt.RegisteredClaims
	UserID string
}

func New(cfg *config.Config, h *handler.Handler, log *slog.Logger) *Transport {
	return &Transport{
		Handler: h,
		Log:     log,
		Config:  cfg,
	}
}

func (t *Transport) NewRouter() *gin.Engine {
	r := gin.Default()

	r.Use(gin.Recovery())
	r.Use(t.withLogging())
	r.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))
	r.Use(t.withCookies())

	r.POST("/api/user/register", func(c *gin.Context) {
		t.Handler.Register(c, *t.Config)
	})

	r.POST("/api/user/login", func(c *gin.Context) {
		t.Handler.Login(c, *t.Config)
	})

	authorized := r.Group("api/user")
	authorized.Use(t.withAuth())

	authorized.POST("/orders", func(c *gin.Context) {
		t.Handler.PostOrders(c, *t.Config)
	})

	authorized.GET("/orders", func(c *gin.Context) {
		t.Handler.GetOrders(c, *t.Config)
	})

	authorized.GET("/balance", func(c *gin.Context) {
		t.Handler.GetBalance(c, *t.Config)
	})

	authorized.POST("/withdraw", func(c *gin.Context) {
		t.Handler.Withdraw(c, *t.Config)
	})

	authorized.GET("/withdrawals", func(c *gin.Context) {
		t.Handler.Withdrawals(c, *t.Config)
	})

	return r
}

func (t *Transport) withLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		uri := c.Request.RequestURI
		method := c.Request.Method

		c.Next()

		status := c.Writer.Status()
		size := c.Writer.Size()
		latency := time.Since(start)

		t.Log.Info("request completed",
			"uri", uri,
			"method", method,
			"duration", latency.String(),
			"status", status,
			"size", size,
		)
		c.Next()
	}
}

func (t *Transport) withCookies() gin.HandlerFunc {
	return func(c *gin.Context) {
		var UserID string
		if cookie, err := c.Cookie("jwt"); err == nil {
			claim := &Claim{}
			token, err := jwt.ParseWithClaims(cookie, claim, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					c.String(http.StatusBadRequest, "Unexpected signing method!")
					return nil, err
				}
				return []byte("123"), nil
			})

			if err != nil {
				c.String(http.StatusUnauthorized, "User ID not found!")
				return
			} else if token.Valid {
				UserID = claim.UserID
				c.Set("user_id", UserID)
				c.Next()
				return
			}
		}

		UserID = uuid.NewString()

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claim{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * time.Minute)),
			},
			UserID: UserID,
		})

		signedToken, err := token.SignedString([]byte("123"))
		if err != nil {
			slog.Error("failed to sign token",
				"error", err,
				"path", c.Request.URL.Path)
			c.Next()
			return
		}

		c.Set("user_id", UserID)
		c.SetCookie("jwt", signedToken, 3600, "/", "", false, true)
		c.Next()
	}
}

func (t *Transport) withAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		UserID := c.GetString("user_id")

		err := t.Handler.Service.Auth(c.Request.Context(), UserID)
		if err != nil {
			if errors.Is(err, storage.ErrUnauthorized) {
				c.AbortWithStatus(401)
				return
			}
		}
		c.Next()
	}
}
