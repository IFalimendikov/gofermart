package transport

import (
	"context"
	"log/slog"
	"time"

	"gofermart/internal/config"
	"gofermart/internal/handler"
	"gofermart/internal/models"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Service interface {
	Register(ctx context.Context, user models.User) error
	Login(ctx context.Context, user models.User) error
	PostOrders(ctx context.Context, login string, orderNum int) error
	GetOrders(ctx context.Context, login string) ([]models.Order, error)
	GetBalance(ctx context.Context, login string) (models.Balance, error)
	Withdraw(ctx context.Context, withdrawal models.Withdrawal) (models.Balance, error)
	Withdrawals(ctx context.Context, login string) ([]models.Withdrawal, error)
}

type Transport struct {
	Handler *handler.Handler
	Log     *slog.Logger
	Config  *config.Config
}

type Claim struct {
	jwt.RegisteredClaims
	Login    string
	Password string
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

	r.POST("/api/user/register", func(c *gin.Context) {
		t.Handler.Register(c)
	})

	r.POST("/api/user/login", func(c *gin.Context) {
		t.Handler.Login(c)
	})

	authorized := r.Group("api/user")
	authorized.Use(t.withCookies())

	authorized.POST("/orders", func(c *gin.Context) {
		t.Handler.PostOrders(c)
	})

	authorized.GET("/orders", func(c *gin.Context) {
		t.Handler.GetOrders(c)
	})

	authorized.GET("/balance", func(c *gin.Context) {
		t.Handler.GetBalance(c)
	})

	authorized.POST("/balance/withdraw", func(c *gin.Context) {
		t.Handler.Withdraw(c)
	})

	authorized.GET("/withdrawals", func(c *gin.Context) {
		t.Handler.Withdrawals(c)
	})

	return r
}

func (t *Transport) withLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now().UTC()
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
		cookie, err := c.Cookie("jwt")
		if err != nil {
			c.AbortWithStatus(401)
			return
		}

		claim := &Claim{}
		token, err := jwt.ParseWithClaims(cookie, claim, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, err
			}
			return []byte("123"), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatus(401)
			return
		}

		c.Set("login", claim.Login)
		c.Set("password", claim.Password)
		c.Next()
	}
}
