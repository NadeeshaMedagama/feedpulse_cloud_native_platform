package handlers

import (
	"encoding/json"
	"net/http"

	"FeedPulse_Cloud_Native_Platform/internal/middleware"
	"FeedPulse_Cloud_Native_Platform/internal/models"
	"FeedPulse_Cloud_Native_Platform/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteJSON(w, http.StatusBadRequest, models.APIResponse{Success: false, Error: "invalid request body"})
		return
	}
	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		middleware.WriteJSON(w, http.StatusUnauthorized, models.APIResponse{Success: false, Error: "invalid email or password"})
		return
	}
	middleware.WriteJSON(w, http.StatusOK, models.APIResponse{Success: true, Data: map[string]string{"token": token}, Message: "login successful"})
}
