package service

import (
	"context"
	"time"

	"go-spring.com/internal/cache"
	"go-spring.com/internal/observability"
	"go-spring.com/internal/repository"
)

// Service handles business logic
type Service struct {
	repository *repository.Repository
	cache      cache.Cache
	tracer     *observability.Tracer
}

// NewService creates a new service instance
func NewService(repo *repository.Repository, cache cache.Cache) *Service {
	return &Service{
		repository: repo,
		cache:      cache,
		tracer:     observability.NewTracer("service"),
	}
}

// GetExample demonstrates service layer business logic with caching
func (s *Service) GetExample() (string, error) {
	ctx := context.Background()
	start := time.Now()

	// Create a cached version of the repository method
	cachedGetExample := cache.Cacheable[string](
		s.cache,
		cache.DefaultKeyGenerator,
		5*time.Minute,
	)(s.repository.GetExample)

	// Call the cached version with tracing
	result, err := observability.TraceFunctionWithResult(s.tracer, ctx, "GetExample", func(ctx context.Context) (string, error) {
		return cachedGetExample()
	})

	// Record metrics
	observability.ServiceMethodDuration.WithLabelValues("Service", "GetExample").Observe(time.Since(start).Seconds())

	return result, err
}

// GetExampleWithParam demonstrates caching with parameters
func (s *Service) GetExampleWithParam(id string) (string, error) {
	ctx := context.Background()
	start := time.Now()

	// Create a cached version of the repository method
	cachedGetExample := cache.Cacheable[string](
		s.cache,
		cache.DefaultKeyGenerator,
		5*time.Minute,
	)(s.repository.GetExample)

	// Call the cached version with tracing
	result, err := observability.TraceFunctionWithResult(s.tracer, ctx, "GetExampleWithParam", func(ctx context.Context) (string, error) {
		return cachedGetExample(id)
	})

	// Record metrics
	observability.ServiceMethodDuration.WithLabelValues("Service", "GetExampleWithParam").Observe(time.Since(start).Seconds())

	return result, err
}
