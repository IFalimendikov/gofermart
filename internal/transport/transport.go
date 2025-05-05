package transport

import (
	"bytes"
	"compress/gzip"
	"io"
	"log/slog"
	"net/http"
	"time"

	"gofermart/internal/config"
	"gofermart/internal/handler"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Transport struct {
	Handler *handler.Handler
	Log *slog.Logger
	Config *config.Config
}

type Claim struct {
	jwt.RegisteredClaims
	UserID string
}

func New(cfg *config.Config, h *handler.Handler, log *slog.Logger) *Transport {
	return &Transport{
		Handler: h,
		Log: log,
		Config: cfg,
	}
}

func(t *Transport) NewRouter() *gin.Engine {
	r := gin.Default()

	r.Use(gin.Recovery())
	r.Use(t.withLogging())
	r.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithDecompressFn(gzip.DefaultDecompressHandle)))
	r.Use(t.withCookies())

	r.POST("/api/user/register", func(c *gin.Context){
		t.Handler.Register(c, *t.Config)
	})

	return r
}

func (t *Transport) withLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		uri := c.Request.RequestURI
		method := c.Request.method

		c.Next()

		status := c.Writer.Status()
		size := c.Writer.Size()
		latency := time.Since(start)

		log.Info("request completed",
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
				if claim.UserID == "" {
					c.String(http.StatusUnauthorized, "User ID not found!")
					return
				}
			} else if token.Valid {
				UserID = claim.UserID
				c.Set("user_id", UserID)
				c.Next()
				return
			}
		}

		UserID = uuid.NewString()

		token := jwt.NewWithClaim(jwt.SigningMethodHS256, Claim{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
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
		c.SetCookie("jwt", signedToken, 60, "/", "", false, true)
		c.Next()
	}
}