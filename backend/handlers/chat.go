package handlers

import (
	"net/http"
	"strings"
	"trading-platform/services"

	"github.com/gin-gonic/gin"
)

type ChatRequest struct {
	Message string `json:"message" binding:"required"`
}

type ChatResponse struct {
	Response string `json:"response"`
}

var tradingKeywords = []string{
	"trade", "trading", "stock", "buy", "sell", "order",
	"portfolio", "account", "balance", "market", "price",
	"deposit", "withdraw", "commission", "fee", "platform",
	"how", "what", "when", "where", "can i",
}

func isTradingRelated(message string) bool {
	messageLower := strings.ToLower(message)
	for _, keyword := range tradingKeywords {
		if strings.Contains(messageLower, keyword) {
			return true
		}
	}
	return false
}

func HandleChat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Filter non-trading questions
	if !isTradingRelated(req.Message) {
		c.JSON(http.StatusOK, ChatResponse{
			Response: "I can only help with questions about our trading platform. Please ask about trading, accounts, orders, deposits, withdrawals, or platform features.",
		})
		return
	}

	// Get response from chatbot service
	response, err := services.GetChatbotResponse(req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ChatResponse{Response: response})
}
