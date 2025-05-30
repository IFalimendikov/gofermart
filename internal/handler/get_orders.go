package handler

import (
	"errors"
	"gofermart/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetOrders(c *gin.Context) {
	login := c.GetString("login")

	orders, err := h.Service.GetOrders(c.Request.Context(), login)
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
