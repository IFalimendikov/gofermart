package handler

import (
	"errors"
	"gofermart/internal/config"
	"gofermart/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetOrders(c *gin.Context, cfg config.Config) {
	userID := c.GetString("login")

	orders, err := h.Service.GetOrders(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, storage.ErrNoOrdersFound) {
			c.JSON(http.StatusNoContent, "")
			return
		} else {
			c.JSON(http.StatusInternalServerError, "")
			return
		}
	}
	c.JSON(http.StatusOK, orders)
}
