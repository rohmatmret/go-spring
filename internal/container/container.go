package container

import (
	"database/sql"
	"fmt"

	"go-spring.com/internal/cache"
	"go-spring.com/internal/config"
	"go-spring.com/internal/repository"
	"go-spring.com/internal/service"
)

type Container struct {
	config   *config.Config
	db       *sql.DB
	cache    cache.Cache
	userRepo *repository.UserRepository
	userSvc  *service.UserService
}

var globalContainer *Container

func NewContainer() (*Container, error) {
	// Load configuration from Vault
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config from Vault: %w", err)
	}

	// Initialize database
	db, err := initDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize cache
	cache, err := initCache(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cache: %w", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	userSvc := service.NewUserService(userRepo, cache)

	container := &Container{
		config:   cfg,
		db:       db,
		cache:    cache,
		userRepo: userRepo,
		userSvc:  userSvc,
	}

	globalContainer = container
	return container, nil
}

func GetContainer() *Container {
	return globalContainer
}

func (c *Container) GetConfig() *config.Config {
	return c.config
}

func (c *Container) GetDB() *sql.DB {
	return c.db
}

func (c *Container) GetCache() cache.Cache {
	return c.cache
}

func (c *Container) GetUserService() *service.UserService {
	return c.userSvc
}

func (c *Container) Close() error {
	if err := c.db.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}
	return nil
}

func initDB(cfg *config.Config) (*sql.DB, error) {

	db, err := sql.Open("postgres", cfg.Database.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

func initCache(cfg *config.Config) (cache.Cache, error) {
	switch cfg.Cache.Type {
	case "memory":
		return cache.NewMemoryCache(), nil
	case "redis":
		return cache.NewRedisCache(
			cfg.Cache.Redis.Host,
			cfg.Cache.Redis.Port,
			cfg.Cache.Redis.Password,
			cfg.Cache.Redis.DB,
		)
	default:
		return nil, fmt.Errorf("unsupported cache type: %s", cfg.Cache.Type)
	}
}
