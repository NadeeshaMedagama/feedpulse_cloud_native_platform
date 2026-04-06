package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"FeedPulse_Cloud_Native_Platform/internal/models"
	"FeedPulse_Cloud_Native_Platform/internal/services"
)

type contextKey string

const AdminEmailContextKey contextKey = "adminEmail"

func JWTAuth(authService *services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				WriteJSON(w, http.StatusUnauthorized, models.APIResponse{Success: false, Error: "missing or invalid token"})
				return
			}
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := authService.Parse(tokenString)
			if err != nil || !token.Valid {
				WriteJSON(w, http.StatusUnauthorized, models.APIResponse{Success: false, Error: "invalid token"})
				return
			}
			ctx := context.WithValue(r.Context(), AdminEmailContextKey, token.Claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func WriteJSON(w http.ResponseWriter, status int, payload models.APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
