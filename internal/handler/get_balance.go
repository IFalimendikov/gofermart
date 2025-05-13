package handler

import (
	"gofermart/internal/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetBalance(c *gin.Context, cfg config.Config) {
	userID := c.GetString("user_id")

	balance, err := h.Service.GetBalance(c.Request.Context(), userID)
	if err != nil {
			c.Status(http.StatusInternalServerError)
			return
	}
	c.JSON(http.StatusAccepted, balance)
}