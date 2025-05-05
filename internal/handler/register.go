package handler

import (
	"errors"
	"fmt"
	"gofermart/internal/config"
	"gofermart/internal/models"
	"gofermart/internal/storage"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Register(c *gin.Context, cfg config.Config) {
	var user models.User

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "Cant read body!")
		return
	}

	if len(body) == 0 {
		c.String(http.StatusBadRequest, "Empty body!")
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		c.String(http.StatusBadRequest, "Error unmarshalling body!")
		return
	}

	urlStr := string(body)
	parsedURL, err := url.Parse(urlStr)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		c.String(http.StatusBadRequest, "Malformed URI!")
		return
	}

	userID := c.GetString("user_id")
	user.ID = userID

	err = h.Service.Register(c.Request.Context(), user)
	if err != nil {
		c.String(http.StatusBadRequest, "Error registering!")
		return
	}

	c.JSON()

}