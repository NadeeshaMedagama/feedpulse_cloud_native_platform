package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"FeedPulse_Cloud_Native_Platform/internal/models"
)

type rateEntry struct {
	count     int
	windowEnd time.Time
}

type IPRateLimiter struct {
	mu      sync.Mutex
	limit   int
	window  time.Duration
	entries map[string]rateEntry
}

func NewIPRateLimiter(limit int, window time.Duration) *IPRateLimiter {
	return &IPRateLimiter{limit: limit, window: window, entries: map[string]rateEntry{}}
}

func (l *IPRateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}
		if ip == "" {
			ip = "unknown"
		}

		now := time.Now()
		l.mu.Lock()
		entry, exists := l.entries[ip]
		if !exists || now.After(entry.windowEnd) {
			entry = rateEntry{count: 0, windowEnd: now.Add(l.window)}
		}
		entry.count++
		l.entries[ip] = entry
		l.mu.Unlock()

		if entry.count > l.limit {
			WriteJSON(w, http.StatusTooManyRequests, models.APIResponse{Success: false, Error: "rate limit exceeded", Message: "Maximum 5 submissions per hour from the same IP."})
			return
		}
		next.ServeHTTP(w, r)
	})
}
