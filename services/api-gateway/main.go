package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	authURL := getEnv("AUTH_SERVICE_URL", "http://localhost:8081")
	feedbackURL := getEnv("FEEDBACK_SERVICE_URL", "http://localhost:8082")
	port := getEnv("GATEWAY_PORT", "8080")

	r := chi.NewRouter()
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(30 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	authProxy := newProxy(authURL)
	feedbackProxy := newProxy(feedbackURL)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/index.html")
	})
	r.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/admin.html")
	})
	r.Handle("/web/*", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))

	r.Handle("/api/auth", authProxy)
	r.Handle("/api/auth/*", authProxy)
	r.Handle("/api/feedback", feedbackProxy)
	r.Handle("/api/feedback/*", feedbackProxy)

	log.Printf("api-gateway running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func newProxy(raw string) *httputil.ReverseProxy {
	u, err := url.Parse(raw)
	if err != nil {
		panic(err)
	}
	return httputil.NewSingleHostReverseProxy(u)
}

func getEnv(k, fallback string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fallback
}
