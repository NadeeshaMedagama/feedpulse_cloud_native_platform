package server

import (
	"net/http"
	"time"

	"FeedPulse_Cloud_Native_Platform/internal/handlers"
	"FeedPulse_Cloud_Native_Platform/internal/middleware"
	"FeedPulse_Cloud_Native_Platform/internal/services"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(authHandler *handlers.AuthHandler, feedbackHandler *handlers.FeedbackHandler, authService *services.AuthService) http.Handler {
	r := chi.NewRouter()
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(60 * time.Second))
	r.Use(cors)

	limiter := middleware.NewIPRateLimiter(5, time.Hour)

	r.Get("/", serveIndex)
	r.Get("/admin", serveAdmin)
	r.Handle("/web/*", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))

	r.Route("/api", func(api chi.Router) {
		api.Post("/auth/login", authHandler.Login)

		api.With(limiter.Middleware).Post("/feedback", feedbackHandler.Create)

		api.Group(func(admin chi.Router) {
			admin.Use(middleware.JWTAuth(authService))
			admin.Get("/feedback", feedbackHandler.List)
			admin.Get("/feedback/summary", feedbackHandler.Summary)
			admin.Get("/feedback/{id}", feedbackHandler.GetByID)
			admin.Patch("/feedback/{id}", feedbackHandler.UpdateStatus)
			admin.Delete("/feedback/{id}", feedbackHandler.Delete)
			admin.Post("/feedback/{id}/reanalyze", feedbackHandler.Reanalyze)
		})
	})

	return r
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/index.html")
}

func serveAdmin(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/admin.html")
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
