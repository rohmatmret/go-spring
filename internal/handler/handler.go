package handler

import (
	"encoding/json"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go-spring.com/internal/container"
	"go-spring.com/internal/service"
)

// Handler manages all HTTP handlers
type Handler struct {
	userHandler *UserHandler
}

// NewHandler creates a new handler
func NewHandler(userService *service.UserService) *Handler {
	return &Handler{
		userHandler: NewUserHandler(userService),
	}
}

// RegisterRoutes registers all routes
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Prometheus metrics endpoint
	mux.Handle("/metrics", promhttp.Handler())

	// Register user routes
	h.userHandler.RegisterRoutes(mux)

	// Example endpoint with caching
	mux.HandleFunc("/api/example", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetExample(w, r, container.GetContainer())
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Example endpoint with parameter and caching
	mux.HandleFunc("/api/example/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetExampleWithParam(w, r, container.GetContainer())
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func handleGetExample(w http.ResponseWriter, r *http.Request, container *container.Container) {
	// Get data from service (will use cache)
	result, err := container.GetUserService().GetUserByID(r.Context(), 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": result,
		"status":  "success",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleGetExampleWithParam(w http.ResponseWriter, r *http.Request, container *container.Container) {
	// Extract ID from URL path
	id := r.URL.Path[len("/api/example/"):]
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	// Get data from service (will use cache)
	result, err := container.GetUserService().GetUserByID(r.Context(), 1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": result,
		"status":  "success",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
