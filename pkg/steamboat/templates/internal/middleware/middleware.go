package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"<<!.ProjectName!>>/internal/middleware/session"
)

func Logger() func(http.Handler) http.Handler {
	return middleware.Logger
}

func CORS() func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	})
}

func Compress() func(http.Handler) http.Handler {
	return middleware.Compress(5)
}

func RequestID() func(http.Handler) http.Handler {
	return middleware.RequestID
}

func Recoverer() func(http.Handler) http.Handler {
	return middleware.Recoverer
}

func RateLimiter(requestsPerMinute int) func(http.Handler) http.Handler {
	return middleware.Throttle(requestsPerMinute)
}

func Timeout(timeout time.Duration) func(http.Handler) http.Handler {
	return middleware.Timeout(timeout)
}

func Session(config *session.Config) func(http.Handler) http.Handler {
	if config == nil {
		config = session.DefaultConfig()
	}
	
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sess, _ := session.LoadSession(r, config)
			ctx := session.WithSession(r.Context(), sess)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}