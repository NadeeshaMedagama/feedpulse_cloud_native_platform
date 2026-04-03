package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"FeedPulse_Cloud_Native_Platform/internal/config"
	"FeedPulse_Cloud_Native_Platform/internal/handlers"
	"FeedPulse_Cloud_Native_Platform/internal/middleware"
	"FeedPulse_Cloud_Native_Platform/internal/repository"
	"FeedPulse_Cloud_Native_Platform/internal/services"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("failed to connect mongodb: %v", err)
	}
	defer func() {
		_ = client.Disconnect(context.Background())
	}()

	feedbackRepo := repository.NewMongoFeedbackRepository(client.Database(cfg.MongoDatabase))
	if err := feedbackRepo.EnsureIndexes(ctx); err != nil {
		log.Printf("index setup warning: %v", err)
	}

	aiServiceURL := os.Getenv("AI_SERVICE_URL")
	if aiServiceURL == "" {
		aiServiceURL = "http://localhost:8083"
	}
	feedbackService := services.NewFeedbackService(feedbackRepo, services.NewHTTPAIClient(aiServiceURL))
	feedbackHandler := handlers.NewFeedbackHandler(feedbackService)
	authService := services.NewAuthService(cfg.AdminEmail, cfg.AdminPassword, cfg.JWTSecret)

	port := os.Getenv("FEEDBACK_SERVICE_PORT")
	if port == "" {
		port = "8082"
	}

	r := chi.NewRouter()
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(60 * time.Second))
	r.Use(cors)

	limiter := middleware.NewIPRateLimiter(5, time.Hour)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	r.With(limiter.Middleware).Post("/api/feedback", feedbackHandler.Create)
	r.Group(func(admin chi.Router) {
		admin.Use(middleware.JWTAuth(authService))
		admin.Get("/api/feedback", feedbackHandler.List)
		admin.Get("/api/feedback/summary", feedbackHandler.Summary)
		admin.Get("/api/feedback/{id}", feedbackHandler.GetByID)
		admin.Patch("/api/feedback/{id}", feedbackHandler.UpdateStatus)
		admin.Delete("/api/feedback/{id}", feedbackHandler.Delete)
		admin.Post("/api/feedback/{id}/reanalyze", feedbackHandler.Reanalyze)
	})

	log.Printf("feedback-service running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
