package handler

import (
	"errors"
	"gofermart/internal/config"
	"gofermart/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetOrders(c *gin.Context, cfg config.Config) {
	userID := c.GetString("user_id")

	orders, err := h.Service.GetOrders(c.Request.Context(), userID)
	if err != nil {
		switch{
		case errors.Is(err, storage.ErrNoOrdersFound):
			c.JSON(http.StatusNoContent, "")
			// c.Status(http.StatusNoContent)
			return
		default:
			c.JSON(http.StatusInternalServerError, "")
			// c.Status(http.StatusBadRequest)
			return
		}
	}
	c.JSON(http.StatusAccepted, orders)
}