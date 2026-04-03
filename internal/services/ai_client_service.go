package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"FeedPulse_Cloud_Native_Platform/internal/models"
)

type HTTPAIClient struct {
	baseURL string
	client  *http.Client
}

func NewHTTPAIClient(baseURL string) *HTTPAIClient {
	return &HTTPAIClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  &http.Client{Timeout: 25 * time.Second},
	}
}

func (c *HTTPAIClient) AnalyzeFeedback(ctx context.Context, title, description string) (models.GeminiAnalysis, error) {
	payload, _ := json.Marshal(models.AIAnalyzeRequest{Title: title, Description: description})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/internal/ai/analyze", bytes.NewReader(payload))
	if err != nil {
		return models.GeminiAnalysis{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return models.GeminiAnalysis{}, err
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return models.GeminiAnalysis{}, fmt.Errorf("ai-service error: %s", string(data))
	}

	var result models.APIResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return models.GeminiAnalysis{}, err
	}
	encoded, _ := json.Marshal(result.Data)
	var analysis models.GeminiAnalysis
	if err := json.Unmarshal(encoded, &analysis); err != nil {
		return models.GeminiAnalysis{}, err
	}
	return analysis, nil
}

func (c *HTTPAIClient) SummarizeThemes(ctx context.Context, text string) (string, error) {
	payload, _ := json.Marshal(models.AISummaryRequest{Text: text})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/internal/ai/summary", bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("ai-service error: %s", string(data))
	}

	var result models.APIResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return "", err
	}
	m, ok := result.Data.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid summary response")
	}
	summary, _ := m["summary"].(string)
	return summary, nil
}
