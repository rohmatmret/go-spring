package repository

import (
	"context"
	"database/sql"
	"time"

	"go-spring.com/internal/observability"
)

// User represents a user entity
type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRepository handles user data access
type UserRepository struct {
	*BaseRepository
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	columns := []string{"id", "username", "email", "created_at", "updated_at"}
	return &UserRepository{
		BaseRepository: NewBaseRepository(db, "users", "id", columns),
	}
}

// FindByUsername finds a user by username
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
	start := time.Now()
	query := "SELECT id, username, email, created_at, updated_at FROM users WHERE username = $1"

	var user User
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	observability.ServiceMethodDuration.WithLabelValues("UserRepository", "FindByUsername").Observe(time.Since(start).Seconds())
	return &user, nil
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	start := time.Now()
	query := "SELECT id, username, email, created_at, updated_at FROM users WHERE email = $1"

	var user User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	observability.ServiceMethodDuration.WithLabelValues("UserRepository", "FindByEmail").Observe(time.Since(start).Seconds())
	return &user, nil
}

// CreateUser creates a new user
func (r *UserRepository) CreateUser(ctx context.Context, user *User) error {
	start := time.Now()
	query := `
		INSERT INTO users (username, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	err := r.db.QueryRowContext(ctx, query,
		user.Username,
		user.Email,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		return err
	}

	observability.ServiceMethodDuration.WithLabelValues("UserRepository", "CreateUser").Observe(time.Since(start).Seconds())
	return nil
}

// UpdateUser updates an existing user
func (r *UserRepository) UpdateUser(ctx context.Context, user *User) error {
	start := time.Now()
	query := `
		UPDATE users
		SET username = $1, email = $2, updated_at = $3
		WHERE id = $4
	`

	user.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, query,
		user.Username,
		user.Email,
		user.UpdatedAt,
		user.ID,
	)

	if err != nil {
		return err
	}

	observability.ServiceMethodDuration.WithLabelValues("UserRepository", "UpdateUser").Observe(time.Since(start).Seconds())
	return nil
}
