package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetBalance(c *gin.Context) {
	login := c.GetString("login")

	balance, err := h.Service.GetBalance(c.Request.Context(), login)
	if err != nil {
		c.JSON(http.StatusNoContent, "")
		return
	}
	c.JSON(http.StatusOK, balance)
}
