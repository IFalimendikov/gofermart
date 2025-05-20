package handler

import (
	"encoding/json"
	"errors"
	"gofermart/internal/config"
	"gofermart/internal/service"
	"gofermart/internal/storage"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) PostOrders(c *gin.Context, cfg config.Config) {
	var orderNum int

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if len(body) == 0 {
		c.Status(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &orderNum)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	login := c.GetString("login")

	err = h.Service.PostOrders(c.Request.Context(), login, orderNum)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrDuplicateOrder):
			c.Status(http.StatusOK)
			return
		case errors.Is(err, storage.ErrDuplicateNumber):
			c.Status(http.StatusConflict)
			return
		case errors.Is(err, service.ErrWrongFormat):
			c.Status(http.StatusUnprocessableEntity)
			return
		default:
			c.Status(http.StatusBadRequest)
			return
		}
	}
	c.Status(http.StatusAccepted)
}
