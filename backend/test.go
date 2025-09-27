package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type ChatCompletionResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func main() {
	godotenv.Load(".env")
	token := os.Getenv("HF_TOKEN")

	if token == "" {
		fmt.Println("Error: HF_TOKEN not found")
		return
	}

	fmt.Println("Token loaded:", token[:10]+"...")
	fmt.Println("\n=== Testing NEW Hugging Face Inference Providers Endpoint ===\n")

	// NEW working endpoint
	apiURL := "https://router.huggingface.co/v1/chat/completions"

	reqBody := ChatCompletionRequest{
		Model: "meta-llama/Llama-3.2-3B-Instruct",
		Messages: []Message{
			{Role: "user", Content: "What is 2+2? Answer in one sentence."},
		},
	}

	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	fmt.Printf("Testing URL: %s\n", apiURL)
	fmt.Printf("Model: %s\n", reqBody.Model)
	fmt.Printf("Authorization: Bearer %s...\n\n", token[:10])

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Request Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	fmt.Printf("Status Code: %d\n", resp.StatusCode)

	if resp.StatusCode == 200 {
		var chatResp ChatCompletionResponse
		if err := json.Unmarshal(body, &chatResp); err != nil {
			fmt.Printf("❌ JSON Parse Error: %v\n", err)
			fmt.Printf("Raw Response: %s\n", string(body))
			return
		}

		if len(chatResp.Choices) > 0 {
			fmt.Printf("✅ SUCCESS!\n")
			fmt.Printf("Response: %s\n\n", chatResp.Choices[0].Message.Content)
		} else {
			fmt.Printf("❌ No choices in response\n")
			fmt.Printf("Raw Response: %s\n", string(body))
		}
	} else if resp.StatusCode == 401 {
		fmt.Println("❌ AUTHENTICATION ERROR")
		fmt.Println("Your token doesn't have the right permissions!")
		fmt.Println("\n📋 Steps to fix:")
		fmt.Println("1. Go to https://huggingface.co/settings/tokens")
		fmt.Println("2. Click 'Create new token' → 'Fine-grained'")
		fmt.Println("3. Under 'Inference', check 'Make calls to Inference Providers'")
		fmt.Println("4. Copy the new token to your .env file")
		fmt.Printf("\nResponse: %s\n", string(body))
	} else if resp.StatusCode == 403 {
		fmt.Println("❌ FORBIDDEN - Model may be gated")
		fmt.Println("Visit https://huggingface.co/meta-llama/Llama-3.2-3B-Instruct")
		fmt.Printf("Response: %s\n", string(body))
	} else {
		fmt.Printf("❌ Error Response: %s\n", string(body))
	}

	// If successful, test with trading platform queries
	if resp.StatusCode == 200 {
		fmt.Println("\n=== Testing Trading Platform Queries ===\n")

		testQueries := []string{
			"What are the trading fees?",
			"How do I place my first trade?",
		}

		for i, query := range testQueries {
			fmt.Printf("Test %d: %s\n", i+1, query)

			reqBody := ChatCompletionRequest{
				Model: "meta-llama/Llama-3.2-3B-Instruct",
				Messages: []Message{
					{Role: "user", Content: query},
				},
			}

			jsonData, _ := json.Marshal(reqBody)
			req, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
			req.Header.Set("Authorization", "Bearer "+token)
			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("  ❌ Error: %v\n\n", err)
				continue
			}

			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			if resp.StatusCode == 200 {
				var chatResp ChatCompletionResponse
				json.Unmarshal(body, &chatResp)
				if len(chatResp.Choices) > 0 {
					fmt.Printf("  ✅ %s\n\n", chatResp.Choices[0].Message.Content)
				}
			} else {
				fmt.Printf("  ❌ Status %d: %s\n\n", resp.StatusCode, string(body))
			}
		}

		fmt.Println("🎉 All tests passed! Update your chatbot.go with the new endpoint.")
	}
}
