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

func (h *Handler) Register(c *gin.Context) {
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
        ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(24 * time.Hour)),
        IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Login:    user.Login,
		Password: user.Password,
	})

	signedToken, err := token.SignedString([]byte(h.cfg.JWTSecret))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.SetCookie("jwt", signedToken, cookieMaxAge, cookiePath, "", false, true)
	c.Status(http.StatusOK)
}
