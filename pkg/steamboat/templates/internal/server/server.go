package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"<<!.ProjectName!>>/internal/database"
	"<<!.ProjectName!>>/internal/handlers"
	"<<!.ProjectName!>>/internal/routes"
	"<<!.ProjectName!>>/internal/utils"
)

type Server struct {
	port     int
	db       database.Service
	handlers *handlers.Handlers
	server   *http.Server
}

func New() *Server {
	return NewWithOptions(true)
}

func NewWithOptions(enableLogging bool) *Server {
	if enableLogging {
		utils.InitLogger()
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil || port == 0 {
		port = 8080
	}

	db := database.New()
	h := handlers.New(db)

	s := &Server{
		port:     port,
		db:       db,
		handlers: h,
	}

	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      routes.Setup(h),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return s
}

func (s *Server) Start() error {
	done := make(chan bool, 1)

	go s.gracefulShutdown(done)

	utils.Logger.Info("Server starting", "port", s.port)

	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("http server error: %w", err)
	}

	<-done
	utils.Logger.Info("Server shutdown complete")

	return nil
}

func (s *Server) gracefulShutdown(done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	utils.Logger.Info("Shutting down gracefully, press Ctrl+C again to force")
	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(shutdownCtx); err != nil {
		utils.Logger.Error("Server forced to shutdown", "error", err)
	}

	if s.db != nil {
		if err := s.db.Close(); err != nil {
			utils.Logger.Error("Error closing database", "error", err)
		} else {
			utils.Logger.Info("Database connection closed")
		}
	}

	utils.Logger.Info("Server exiting")

	done <- true
}

