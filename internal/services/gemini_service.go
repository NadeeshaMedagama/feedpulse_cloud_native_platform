package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"FeedPulse_Cloud_Native_Platform/internal/models"
)

type GeminiService interface {
	AnalyzeFeedback(ctx context.Context, title, description string) (models.GeminiAnalysis, error)
	SummarizeThemes(ctx context.Context, text string) (string, error)
}

type HTTPGeminiService struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewGeminiService(apiKey string) *HTTPGeminiService {
	return &HTTPGeminiService{
		apiKey:  apiKey,
		baseURL: "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent",
		client:  &http.Client{},
	}
}

func (s *HTTPGeminiService) AnalyzeFeedback(ctx context.Context, title, description string) (models.GeminiAnalysis, error) {
	if s.apiKey == "" {
		return models.GeminiAnalysis{}, errors.New("gemini api key not set")
	}

	prompt := fmt.Sprintf("Analyse this product feedback. Return ONLY valid JSON with fields: category, sentiment, priority_score (1-10), summary, tags.\nTitle: %s\nDescription: %s", title, description)
	text, err := s.call(ctx, prompt)
	if err != nil {
		return models.GeminiAnalysis{}, err
	}

	payload := extractJSON(text)
	var analysis models.GeminiAnalysis
	if err := json.Unmarshal([]byte(payload), &analysis); err != nil {
		return models.GeminiAnalysis{}, err
	}
	return analysis, nil
}

func (s *HTTPGeminiService) SummarizeThemes(ctx context.Context, text string) (string, error) {
	if s.apiKey == "" {
		return "Gemini API key not configured.", errors.New("gemini api key not set")
	}
	prompt := "Summarize top 3 themes from this last 7 days feedback. Keep it concise in bullet style.\n" + text
	return s.call(ctx, prompt)
}

func (s *HTTPGeminiService) call(ctx context.Context, prompt string) (string, error) {
	url := fmt.Sprintf("%s?key=%s", strings.TrimRight(s.baseURL, "?"), s.apiKey)
	bodyMap := map[string]interface{}{
		"contents": []map[string]interface{}{{
			"parts": []map[string]string{{"text": prompt}},
		}},
	}
	body, _ := json.Marshal(bodyMap)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("gemini request failed: %s", string(data))
	}

	var parsed struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return "", err
	}
	if len(parsed.Candidates) == 0 || len(parsed.Candidates[0].Content.Parts) == 0 {
		return "", errors.New("empty gemini response")
	}
	return parsed.Candidates[0].Content.Parts[0].Text, nil
}

func extractJSON(s string) string {
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start == -1 || end == -1 || end < start {
		return s
	}
	return s[start : end+1]
}
