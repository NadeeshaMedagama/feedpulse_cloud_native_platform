package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"FeedPulse_Cloud_Native_Platform/internal/config"
	"FeedPulse_Cloud_Native_Platform/internal/handlers"
	"FeedPulse_Cloud_Native_Platform/internal/services"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.Load()
	authService := services.NewAuthService(cfg.AdminEmail, cfg.AdminPassword, cfg.JWTSecret)
	authHandler := handlers.NewAuthHandler(authService)

	port := os.Getenv("AUTH_SERVICE_PORT")
	if port == "" {
		port = "8081"
	}

	r := chi.NewRouter()
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(20 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	r.Post("/api/auth/login", authHandler.Login)

	log.Printf("auth-service running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
