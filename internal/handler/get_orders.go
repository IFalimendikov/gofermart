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
        if errors.Is(err, storage.ErrNoOrdersFound) {
            c.Status(http.StatusNoContent) // 204 если нет заказов
            return
        }
        c.Status(http.StatusInternalServerError) // 500 для других ошибок
        return
    }
	    c.Header("Content-Type", "application/json")
    c.JSON(http.StatusOK, orders)
}