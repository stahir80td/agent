package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const platformDocs = `
# Trading Platform Documentation

## Getting Started
- Sign up for a free account at tradepro.com
- Complete KYC verification within 24 hours
- Fund your account via bank transfer (minimum $100)

## Trading Features
- Buy and sell stocks in real-time during market hours (9:30 AM - 4:00 PM EST)
- Place market orders (instant execution) and limit orders (set your price)
- View live charts and market data with real-time updates
- Set stop-loss orders to limit losses and take-profit orders to secure gains
- Trade options with $0.65 per contract fee

## Account Management
- Deposit funds: Minimum $100, instant via debit card or 3-5 days for ACH
- Withdraw funds: 2-3 business days processing time, $5 withdrawal fee
- View complete transaction history and detailed statements
- Check portfolio performance with real-time P&L tracking
- Set up recurring deposits for dollar-cost averaging

## Fees and Commissions
- Stock trades: $0 commission (completely free!)
- Options trades: $0.65 per contract
- Cryptocurrency: 1% spread on all crypto trades
- Withdrawal fee: $5 per withdrawal
- No account maintenance fees, no minimum balance required

## Common Questions

Q: How do I place my first trade?
A: Click "Trade" in the top menu → Search for stock → Choose order type (market/limit) → Enter quantity → Review and confirm

Q: What are the trading hours?
A: Regular market hours are 9:30 AM - 4:00 PM EST, Monday-Friday. Pre-market: 4:00 AM - 9:30 AM, After-hours: 4:00 PM - 8:00 PM

Q: How long do deposits take?
A: Debit card deposits are instant. ACH bank transfers take 3-5 business days to clear

Q: Is my money safe?
A: Yes! We're SIPC insured up to $500,000, and we use bank-level encryption

Q: Can I trade on mobile?
A: Yes! Download our iOS or Android app for full trading capabilities on the go
`

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
}

type ChatCompletionResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func GetChatbotResponse(query string) (string, error) {
	hfToken := os.Getenv("HF_TOKEN")
	if hfToken == "" {
		return "Configuration error: HF_TOKEN not set", nil
	}

	// Build messages with system context
	systemMessage := fmt.Sprintf(`You are a helpful trading platform assistant. Answer questions using ONLY the information provided in the documentation below. If the answer is not in the documentation, say "I don't have that information in my knowledge base."

Documentation:
%s`, platformDocs)

	reqBody := ChatCompletionRequest{
		Model: "meta-llama/Llama-3.2-3B-Instruct",
		Messages: []Message{
			{Role: "system", Content: systemMessage},
			{Role: "user", Content: query},
		},
		MaxTokens:   200,
		Temperature: 0.7,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	// NEW WORKING ENDPOINT
	apiURL := "https://router.huggingface.co/v1/chat/completions"
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+hfToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Check for errors
	if resp.StatusCode != 200 {
		return fmt.Sprintf("API Error (Status %d): %s", resp.StatusCode, string(body)), nil
	}

	// Parse chat completion response
	var chatResp ChatCompletionResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", err
	}

	if len(chatResp.Choices) == 0 {
		return "I couldn't generate a response. Please try again.", nil
	}

	response := strings.TrimSpace(chatResp.Choices[0].Message.Content)
	return response, nil
}
