package handler

import (
	"encoding/json"
	"errors"
	"gofermart/internal/models"
	"gofermart/internal/storage"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	tokenExpiryHours = 24
	cookieMaxAge     = 24 * 3600
	cookiePath       = "/"
)

type Claim struct {
	jwt.RegisteredClaims
	Login    string
	Password string
}

func (h *Handler) Login(c *gin.Context) {
	var user models.User

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if len(body) == 0 {
		c.Status(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = h.Service.Login(c.Request.Context(), user)
	if err != nil {
		if errors.Is(err, storage.ErrWrongPassword) {
			c.Status(http.StatusUnauthorized)
			return
		}
		c.Status(http.StatusBadRequest)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claim{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Login:    user.Login,
		Password: user.Password,
	})

	signedToken, err := token.SignedString([]byte("123"))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.SetCookie("jwt", signedToken, cookieMaxAge, cookiePath, "", false, true)
	c.Status(http.StatusOK)
}
