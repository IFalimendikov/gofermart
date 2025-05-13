package handler

import (
	"errors"
	"gofermart/internal/config"
	"gofermart/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Withdrawals(c *gin.Context, cfg config.Config) {
	userID := c.GetString("user_id")
	
	withdrawals, err := h.Service.Withdrawals(c.Request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrNoOrdersFound):
			c.Status(http.StatusNoContent)
			return
		default:
			c.Status(http.StatusBadRequest)
			return
		}
	}
	c.JSON(http.StatusOK, withdrawals)
}
