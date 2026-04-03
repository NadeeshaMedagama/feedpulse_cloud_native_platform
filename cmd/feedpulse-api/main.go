package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"FeedPulse_Cloud_Native_Platform/internal/config"
	"FeedPulse_Cloud_Native_Platform/internal/handlers"
	"FeedPulse_Cloud_Native_Platform/internal/repository"
	"FeedPulse_Cloud_Native_Platform/internal/server"
	"FeedPulse_Cloud_Native_Platform/internal/services"

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

	db := client.Database(cfg.MongoDatabase)

	feedbackRepo := repository.NewMongoFeedbackRepository(db)
	if err := feedbackRepo.EnsureIndexes(ctx); err != nil {
		log.Printf("index setup warning: %v", err)
	}

	geminiService := services.NewGeminiService(cfg.GeminiAPIKey)
	feedbackService := services.NewFeedbackService(feedbackRepo, geminiService)
	authService := services.NewAuthService(cfg.AdminEmail, cfg.AdminPassword, cfg.JWTSecret)

	authHandler := handlers.NewAuthHandler(authService)
	feedbackHandler := handlers.NewFeedbackHandler(feedbackService)

	router := server.NewRouter(authHandler, feedbackHandler, authService)

	addr := ":" + cfg.Port
	log.Printf("FeedPulse API running on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
