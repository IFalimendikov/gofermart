package handler

import (
	"encoding/json"
	"errors"
	"gofermart/internal/config"
	"gofermart/internal/models"
	"gofermart/internal/service"
	"gofermart/internal/storage"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)
func (h *Handler) Withdraw(c *gin.Context, cfg config.Config) {
	var withdrawal models.Withdrawal
	var balance models.Balance

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		c.Status(http.StatusBadRequest)
		return
	}

	if len(body) == 0 {
		log.Print("Empty request body received")
		c.Status(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &withdrawal)
	if err != nil {
		log.Printf("Error unmarshaling withdrawal request: %v", err)
		c.Status(http.StatusBadRequest)
		return
	}

	userID := c.GetString("login")
	withdrawal.ID = userID
	
	log.Printf("Processing withdrawal request - UserID: %s, Amount: %v, Order: %s", 
		userID, withdrawal.Sum, withdrawal.Order)
	
	balance, err = h.Service.Withdraw(c.Request.Context(), withdrawal)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrNoOrdersFound):
			log.Printf("No orders found for withdrawal - UserID: %s", userID)
			c.Status(http.StatusNoContent)
			return
		case errors.Is(err, service.ErrWrongFormat):
			log.Printf("Wrong format in withdrawal request: %v", err)
			c.Status(http.StatusUnprocessableEntity)
			return
		default:
			log.Printf("Failed to process withdrawal: %v", err)
			c.Status(http.StatusBadRequest)
			return
		}
	}

	log.Printf("Withdrawal successful - UserID: %s, New Balance: %v", 
		userID, balance.Current)
	
	c.JSON(http.StatusOK, balance)
}