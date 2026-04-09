package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"FeedPulse_Cloud_Native_Platform/internal/middleware"
	"FeedPulse_Cloud_Native_Platform/internal/models"
	"FeedPulse_Cloud_Native_Platform/internal/services"

	"github.com/go-chi/chi/v5"
)

type FeedbackHandler struct {
	feedbackService *services.FeedbackService
}

func NewFeedbackHandler(feedbackService *services.FeedbackService) *FeedbackHandler {
	return &FeedbackHandler{feedbackService: feedbackService}
}

func (h *FeedbackHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.FeedbackCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteJSON(w, http.StatusBadRequest, models.APIResponse{Success: false, Error: "invalid request body"})
		return
	}
	feedback, err := h.feedbackService.Create(r.Context(), req)
	if err != nil {
		middleware.WriteJSON(w, http.StatusBadRequest, models.APIResponse{Success: false, Error: err.Error()})
		return
	}
	middleware.WriteJSON(w, http.StatusCreated, models.APIResponse{Success: true, Data: feedback, Message: "feedback submitted"})
}

func (h *FeedbackHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	filters := models.FeedbackFilters{
		Category: r.URL.Query().Get("category"),
		Status:   r.URL.Query().Get("status"),
		Search:   r.URL.Query().Get("search"),
		SortBy:   r.URL.Query().Get("sortBy"),
		Order:    r.URL.Query().Get("order"),
		Page:     page,
		Limit:    limit,
	}
	items, total, err := h.feedbackService.List(r.Context(), filters)
	if err != nil {
		middleware.WriteJSON(w, http.StatusInternalServerError, models.APIResponse{Success: false, Error: "failed to fetch feedback"})
		return
	}
	stats, err := h.feedbackService.Stats(r.Context())
	if err != nil {
		stats = map[string]interface{}{}
	}
	middleware.WriteJSON(w, http.StatusOK, models.APIResponse{Success: true, Data: map[string]interface{}{"items": items, "total": total, "page": filters.Page, "limit": filters.Limit, "stats": stats}})
}

func (h *FeedbackHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	item, err := h.feedbackService.GetByID(r.Context(), id)
	if err != nil {
		status := http.StatusBadRequest
		if services.IsNotFound(err) {
			status = http.StatusNotFound
		}
		middleware.WriteJSON(w, status, models.APIResponse{Success: false, Error: "feedback not found"})
		return
	}
	middleware.WriteJSON(w, http.StatusOK, models.APIResponse{Success: true, Data: item})
}

func (h *FeedbackHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req models.FeedbackUpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteJSON(w, http.StatusBadRequest, models.APIResponse{Success: false, Error: "invalid request body"})
		return
	}
	if err := h.feedbackService.UpdateStatus(r.Context(), id, req.Status); err != nil {
		status := http.StatusBadRequest
		if services.IsNotFound(err) {
			status = http.StatusNotFound
		}
		middleware.WriteJSON(w, status, models.APIResponse{Success: false, Error: err.Error()})
		return
	}
	middleware.WriteJSON(w, http.StatusOK, models.APIResponse{Success: true, Message: "status updated"})
}

func (h *FeedbackHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.feedbackService.Delete(r.Context(), id); err != nil {
		status := http.StatusBadRequest
		if services.IsNotFound(err) {
			status = http.StatusNotFound
		}
		middleware.WriteJSON(w, status, models.APIResponse{Success: false, Error: "failed to delete feedback"})
		return
	}
	middleware.WriteJSON(w, http.StatusOK, models.APIResponse{Success: true, Message: "feedback deleted"})
}

func (h *FeedbackHandler) Summary(w http.ResponseWriter, r *http.Request) {
	summary, err := h.feedbackService.WeeklySummary(r.Context())
	if err != nil {
		middleware.WriteJSON(w, http.StatusInternalServerError, models.APIResponse{Success: false, Error: "failed to generate summary"})
		return
	}
	middleware.WriteJSON(w, http.StatusOK, models.APIResponse{Success: true, Data: map[string]string{"summary": summary}})
}

func (h *FeedbackHandler) Reanalyze(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.feedbackService.Reanalyze(r.Context(), id); err != nil {
		middleware.WriteJSON(w, http.StatusBadRequest, models.APIResponse{Success: false, Error: err.Error()})
		return
	}
	middleware.WriteJSON(w, http.StatusOK, models.APIResponse{Success: true, Message: "AI analysis retriggered"})
}
