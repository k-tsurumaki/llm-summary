package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// LLMClient manages communication with the LLM API
type LLMClient struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

type LLMRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type LLMResponse struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Choices []struct {
		Index        int     `json:"index"`
		FinishReason string  `json:"finish_reason"`
		Message      Message `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// NewLLMClient creates a new LLMClient instance
func NewLLMClient() *LLMClient {
	return &LLMClient{
		apiKey:  Config.LLM_API_KEY,
		baseURL: Config.LLM_BASE_URL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Summarize sends a prompt to the LLM and returns the summarized content
func (c *LLMClient) Summarize(ctx context.Context, text string) (string, error) {
	if c.apiKey == "" {
		return "", fmt.Errorf("LLM API key not configured")
	}

	prompt := fmt.Sprintf("以下の文章を簡潔に要約してください：\n\n%s", text)
	log.Printf("Prompt for LLM: %s", prompt)

	reqBody := LLMRequest{
		Model: "google/gemma-3-27b-it:free", // 修正: モデル名が正しいことを確認
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	log.Printf("LLM raw response: %s", string(bodyBytes))

	var llmResp LLMResponse
	if err := json.Unmarshal(bodyBytes, &llmResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(llmResp.Choices) == 0 || llmResp.Choices[0].Message.Content == "" {
		return "", fmt.Errorf("no valid response from LLM: %+v", llmResp)
	}

	return llmResp.Choices[0].Message.Content, nil
}
