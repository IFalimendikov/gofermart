package handler

import (
	"encoding/json"
	"errors"
	"gofermart/internal/config"
	"gofermart/internal/models"
	"gofermart/internal/storage"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (h *Handler) Register(c *gin.Context, cfg config.Config) {
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

	err = h.Service.Register(c.Request.Context(), user)
	if err != nil {
		if errors.Is(err, storage.ErrDuplicateLogin) {
			c.Status(http.StatusConflict)
			return
		}
		c.Status(http.StatusBadRequest)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claim{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Login:    user.Login,
		Password: user.Password,
	})

	signedToken, err := token.SignedString([]byte("123"))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.SetCookie("jwt", signedToken, 24*3600, "/", "", false, true)
	c.Status(http.StatusOK)
}
