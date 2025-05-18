package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"go-spring.com/internal/observability"
	"go-spring.com/internal/repository"
	"go-spring.com/internal/service"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService *service.UserService
	tracer      *observability.Tracer
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		tracer:      observability.NewTracer("user_handler"),
	}
}

// RegisterRoutes registers the user routes
func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/users", h.handleUsers)
	mux.HandleFunc("/api/users/", h.handleUser)
}

// handleUsers handles /api/users endpoints
func (h *UserHandler) handleUsers(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		h.createUser(w, r)
	case http.MethodGet:
		h.getUserByUsername(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleUser handles /api/users/{id} endpoints
func (h *UserHandler) handleUser(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from URL
	idStr := r.URL.Path[len("/api/users/"):]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getUserByID(w, r, id)
	case http.MethodPut:
		h.updateUser(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// createUser handles POST /api/users
func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse request body
	var user repository.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create user
	err := observability.TraceFunction(h.tracer, ctx, "CreateUser", func(ctx context.Context) error {
		return h.userService.CreateUser(ctx, &user)
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return created user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// getUserByID handles GET /api/users/{id}
func (h *UserHandler) getUserByID(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := r.Context()

	// Get user
	user, err := observability.TraceFunctionWithResult(h.tracer, ctx, "GetUserByID", func(ctx context.Context) (*repository.User, error) {
		return h.userService.GetUserByID(ctx, id)
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Return user
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// getUserByUsername handles GET /api/users?username={username}
func (h *UserHandler) getUserByUsername(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get username from query parameter
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	// Get user
	user, err := observability.TraceFunctionWithResult(h.tracer, ctx, "GetUserByUsername", func(ctx context.Context) (*repository.User, error) {
		return h.userService.GetUserByUsername(ctx, username)
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Return user
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// updateUser handles PUT /api/users/{id}
func (h *UserHandler) updateUser(w http.ResponseWriter, r *http.Request, id int64) {
	ctx := r.Context()

	// Parse request body
	var user repository.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set user ID from URL
	user.ID = id

	// Update user
	err := observability.TraceFunction(h.tracer, ctx, "UpdateUser", func(ctx context.Context) error {
		return h.userService.UpdateUser(ctx, &user)
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return updated user
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
