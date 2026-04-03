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
	gemini := services.NewGeminiService(cfg.GeminiAPIKey)
	handler := handlers.NewAIHandler(gemini)

	port := os.Getenv("AI_SERVICE_PORT")
	if port == "" {
		port = "8083"
	}

	r := chi.NewRouter()
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(30 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	r.Post("/internal/ai/analyze", handler.Analyze)
	r.Post("/internal/ai/summary", handler.Summary)

	log.Printf("ai-service running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
