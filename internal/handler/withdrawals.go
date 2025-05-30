package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Withdrawals(c *gin.Context) {
	login := c.GetString("login")

	withdrawals, err := h.Service.Withdrawals(c.Request.Context(), login)
	if err != nil {
		switch {
		default:
			c.Status(http.StatusBadRequest)
			return
		}
	}

	if len(withdrawals) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, withdrawals)
}
