package main

import (
	"net/http"

	"github.com/k-tsurumaki/fuselage"
)

type SummaryRequest struct {
	Text string `json:"text"`
}

type SummaryResponse struct {
	Summary string `json:"summary"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func setupRoutes(app *fuselage.Router) {

	// ヘルスチェック
	app.GET("/health", healthHandler)

	// 要約API
	app.POST("/summarize", summarizeHandler)
}

func healthHandler(c *fuselage.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func summarizeHandler(c *fuselage.Context) error {
	var req SummaryRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid JSON format"})
	}

	if req.Text == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Text field is required"})
	}

	// LLMクライアントを使用して要約を取得
	llmClient := NewLLMClient()
	summary, err := llmClient.Summarize(c.Request.Context(), req.Text)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to generate summary"})
	}

	return c.JSON(http.StatusOK, SummaryResponse{Summary: summary})
}
