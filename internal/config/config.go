package config

import (
	"os"
)

type Config struct {
	Port          string
	MongoURI      string
	MongoDatabase string
	JWTSecret     string
	AdminEmail    string
	AdminPassword string
	GeminiAPIKey  string
}

func Load() Config {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = os.Getenv("MONGO_LOCAL_URI")
	}
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	cfg := Config{
		Port:          getEnv("PORT", "8080"),
		MongoURI:      mongoURI,
		MongoDatabase: getEnv("MONGO_DATABASE", "feedpulse"),
		JWTSecret:     getEnv("JWT_SECRET", "dev-secret-change-me"),
		AdminEmail:    getEnv("ADMIN_EMAIL", "admin@feedpulse.local"),
		AdminPassword: getEnv("ADMIN_PASSWORD", "Admin123!"),
		GeminiAPIKey:  os.Getenv("GEMINI_API_KEY"),
	}
	return cfg
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
