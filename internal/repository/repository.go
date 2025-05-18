package repository

import (
	"fmt"

	"go-spring.com/internal/config"
)

// Repository handles data access
type Repository struct {
	config *config.Config
}

// NewRepository creates a new repository instance
func NewRepository(cfg *config.Config) *Repository {
	return &Repository{
		config: cfg,
	}
}

// GetExample demonstrates repository pattern
func (r *Repository) GetExample(args ...interface{}) (string, error) {
	// This is where you would typically interact with a database
	if len(args) > 0 {
		if id, ok := args[0].(string); ok {
			return fmt.Sprintf("Example data for ID: %s", id), nil
		}
	}
	return "Example data from repository", nil
}
