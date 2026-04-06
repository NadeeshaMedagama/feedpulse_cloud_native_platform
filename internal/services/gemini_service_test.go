package services

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGeminiServiceAnalyzeFeedback_ParsesJSONPayload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
  "candidates": [
    {
      "content": {
        "parts": [
          {
					  "text": "Result: {\"category\":\"Feature Request\",\"sentiment\":\"Positive\",\"priority_score\":8,\"summary\":\"User wants dark mode.\",\"tags\":[\"UI\",\"Settings\"]}"
          }
        ]
      }
    }
  ]
}`))
	}))
	defer server.Close()

	svc := NewGeminiService("test-key")
	svc.baseURL = server.URL
	svc.client = server.Client()

	analysis, err := svc.AnalyzeFeedback(context.Background(), "Need dark mode", "Please add dark mode in dashboard settings.")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if analysis.Category != "Feature Request" {
		t.Fatalf("expected category Feature Request, got %q", analysis.Category)
	}
	if analysis.PriorityScore != 8 {
		t.Fatalf("expected priority 8, got %d", analysis.PriorityScore)
	}
	if len(analysis.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(analysis.Tags))
	}
}
