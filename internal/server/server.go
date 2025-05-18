package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go-spring.com/internal/container"
	"go-spring.com/internal/handler"
	"go-spring.com/internal/observability"
)

// Server represents the HTTP server
type Server struct {
	container *container.Container
	server    *http.Server
	tracer    *observability.Tracer
}

// NewServer creates a new HTTP server
func NewServer(container *container.Container) *Server {
	cfg := container.GetConfig()

	// Create router and register routes
	router := http.NewServeMux()

	// Initialize tracer
	tracer := observability.NewTracer("go-spring")

	// Add observability middleware
	router = http.NewServeMux()
	handler.NewHandler(container.GetUserService()).RegisterRoutes(router)

	// Create HTTP server with middleware
	server := &http.Server{
		Addr: fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler: observability.MetricsMiddleware(
			observability.TracingMiddleware(tracer)(router),
		),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		container: container,
		server:    server,
		tracer:    tracer,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
