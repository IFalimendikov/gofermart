package handler

import (
	"encoding/json"
	"errors"
	"gofermart/internal/config"
	"gofermart/internal/models"
	"gofermart/internal/storage"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
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

	userID := c.GetString("user_id")
	user.ID = userID

	err = h.Service.Register(c.Request.Context(), user)
	if err != nil {
		if errors.Is(err, storage.ErrDuplicateLogin) {
			c.Status(http.StatusConflict)
			return
		}
		c.Status(http.StatusBadRequest)
		return
	}

	c.Status(http.StatusOK)
}