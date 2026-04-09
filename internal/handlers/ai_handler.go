package handlers

import (
	"encoding/json"
	"net/http"

	"FeedPulse_Cloud_Native_Platform/internal/middleware"
	"FeedPulse_Cloud_Native_Platform/internal/models"
	"FeedPulse_Cloud_Native_Platform/internal/services"
)

type AIHandler struct {
	gemini services.GeminiService
}

func NewAIHandler(gemini services.GeminiService) *AIHandler {
	return &AIHandler{gemini: gemini}
}

func (h *AIHandler) Analyze(w http.ResponseWriter, r *http.Request) {
	var req models.AIAnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteJSON(w, http.StatusBadRequest, models.APIResponse{Success: false, Error: "invalid request body"})
		return
	}
	analysis, err := h.gemini.AnalyzeFeedback(r.Context(), req.Title, req.Description)
	if err != nil {
		middleware.WriteJSON(w, http.StatusBadGateway, models.APIResponse{Success: false, Error: err.Error()})
		return
	}
	middleware.WriteJSON(w, http.StatusOK, models.APIResponse{Success: true, Data: analysis})
}

func (h *AIHandler) Summary(w http.ResponseWriter, r *http.Request) {
	var req models.AISummaryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteJSON(w, http.StatusBadRequest, models.APIResponse{Success: false, Error: "invalid request body"})
		return
	}
	summary, err := h.gemini.SummarizeThemes(r.Context(), req.Text)
	if err != nil {
		middleware.WriteJSON(w, http.StatusBadGateway, models.APIResponse{Success: false, Error: err.Error()})
		return
	}
	middleware.WriteJSON(w, http.StatusOK, models.APIResponse{Success: true, Data: map[string]string{"summary": summary}})
}
