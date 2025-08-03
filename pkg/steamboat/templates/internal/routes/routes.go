package routes

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"<<!.ProjectName!>>/internal/handlers"
	"<<!.ProjectName!>>/internal/middleware"
)

func Setup(h *handlers.Handlers) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer())
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(middleware.Session(nil))
	r.Use(middleware.RateLimiter(100))
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.Compress())
	r.Use(middleware.CORS())

	//User routes
	r.Get("/", h.HomeHandler)

	return r
}

