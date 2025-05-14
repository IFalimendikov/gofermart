package handler

import (
	"encoding/json"
	"errors"
	"gofermart/internal/config"
	"gofermart/internal/models"
	"gofermart/internal/service"
	"gofermart/internal/storage"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Withdraw(c *gin.Context, cfg config.Config) {
	var withdrawal models.Withdrawal
	var balance models.Balance

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if len(body) == 0 {
		c.Status(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &withdrawal)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	userID := c.GetString("login")
	withdrawal.ID = userID
	
	balance, err = h.Service.Withdraw(c.Request.Context(), withdrawal)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrNoOrdersFound):
			c.Status(http.StatusNoContent)
			return
		case errors.Is(err, service.ErrWrongFormat):
			c.Status(http.StatusUnprocessableEntity)
			return
		default:
			c.Status(http.StatusBadRequest)
			return
		}
	}
	c.JSON(http.StatusOK, balance)
}
