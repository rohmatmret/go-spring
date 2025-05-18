package service

import (
	"context"
	"fmt"
	"time"

	"go-spring.com/internal/cache"
	"go-spring.com/internal/observability"
	"go-spring.com/internal/repository"
)

// UserService handles user-related business logic
type UserService struct {
	userRepo *repository.UserRepository
	cache    cache.Cache
	tracer   *observability.Tracer
}

// NewUserService creates a new user service
func NewUserService(userRepo *repository.UserRepository, cache cache.Cache) *UserService {
	return &UserService{
		userRepo: userRepo,
		cache:    cache,
		tracer:   observability.NewTracer("user_service"),
	}
}

func (s *UserService) GetUsers(ctx context.Context) ([]*repository.User, error) {
	return nil, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(ctx context.Context, id int64) (*repository.User, error) {
	start := time.Now()

	// Try to get from cache first
	cacheKey := fmt.Sprintf("user:%d", id)
	if cached, err := s.cache.Get(ctx, cacheKey); err {
		if user, ok := cached.(*repository.User); ok {
			observability.CacheHits.WithLabelValues().Inc()
			observability.ServiceMethodDuration.WithLabelValues("UserService", "GetUserByID").Observe(time.Since(start).Seconds())
			return user, nil
		}
	}
	observability.CacheMisses.WithLabelValues().Inc()

	// Get from database
	user, err := observability.TraceFunctionWithResult(s.tracer, ctx, "GetUserByID", func(ctx context.Context) (*repository.User, error) {
		return s.userRepo.FindByID(ctx, id)
	})
	if err != nil {
		return nil, err
	}

	// Cache the result
	if user != nil {
		s.cache.Set(ctx, cacheKey, user, 5*time.Minute)
	}

	observability.ServiceMethodDuration.WithLabelValues("UserService", "GetUserByID").Observe(time.Since(start).Seconds())
	return user, nil
}

// GetUserByUsername retrieves a user by username
func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*repository.User, error) {
	start := time.Now()

	// Try to get from cache first
	cacheKey := fmt.Sprintf("user:username:%s", username)
	if cached, err := s.cache.Get(ctx, cacheKey); err {
		if user, ok := cached.(*repository.User); ok {
			observability.CacheHits.WithLabelValues().Inc()
			observability.ServiceMethodDuration.WithLabelValues("UserService", "GetUserByUsername").Observe(time.Since(start).Seconds())
			return user, nil
		}
	}
	observability.CacheMisses.WithLabelValues().Inc()

	// Get from database
	user, err := observability.TraceFunctionWithResult(s.tracer, ctx, "GetUserByUsername", func(ctx context.Context) (*repository.User, error) {
		return s.userRepo.FindByUsername(ctx, username)
	})
	if err != nil {
		return nil, err
	}

	// Cache the result
	if user != nil {
		s.cache.Set(ctx, cacheKey, user, 5*time.Minute)
	}

	observability.ServiceMethodDuration.WithLabelValues("UserService", "GetUserByUsername").Observe(time.Since(start).Seconds())
	return user, nil
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, user *repository.User) error {
	start := time.Now()

	// Check if username exists
	existingUser, err := s.userRepo.FindByUsername(ctx, user.Username)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return fmt.Errorf("username already exists")
	}

	// Check if email exists
	existingUser, err = s.userRepo.FindByEmail(ctx, user.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return fmt.Errorf("email already exists")
	}

	// Create user
	err = observability.TraceFunction(s.tracer, ctx, "CreateUser", func(ctx context.Context) error {
		return s.userRepo.CreateUser(ctx, user)
	})
	if err != nil {
		return err
	}

	// Invalidate cache
	s.cache.Delete(ctx, fmt.Sprintf("user:%d", user.ID))
	s.cache.Delete(ctx, fmt.Sprintf("user:username:%s", user.Username))

	observability.ServiceMethodDuration.WithLabelValues("UserService", "CreateUser").Observe(time.Since(start).Seconds())
	return nil
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(ctx context.Context, user *repository.User) error {
	start := time.Now()

	// Check if user exists
	existingUser, err := s.userRepo.FindByID(ctx, user.ID)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return fmt.Errorf("user not found")
	}

	// Check if new username is taken
	if user.Username != existingUser.Username {
		existingUser, err = s.userRepo.FindByUsername(ctx, user.Username)
		if err != nil {
			return err
		}
		if existingUser != nil {
			return fmt.Errorf("username already exists")
		}
	}

	// Check if new email is taken
	if user.Email != existingUser.Email {
		existingUser, err = s.userRepo.FindByEmail(ctx, user.Email)
		if err != nil {
			return err
		}
		if existingUser != nil {
			return fmt.Errorf("email already exists")
		}
	}

	// Update user
	err = observability.TraceFunction(s.tracer, ctx, "UpdateUser", func(ctx context.Context) error {
		return s.userRepo.UpdateUser(ctx, user)
	})
	if err != nil {
		return err
	}

	// Invalidate cache
	s.cache.Delete(ctx, fmt.Sprintf("user:%d", user.ID))
	s.cache.Delete(ctx, fmt.Sprintf("user:username:%s", user.Username))

	observability.ServiceMethodDuration.WithLabelValues("UserService", "UpdateUser").Observe(time.Since(start).Seconds())
	return nil
}
