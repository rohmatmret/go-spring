package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-spring.com/internal/config"
	"go-spring.com/internal/container"
	"go-spring.com/internal/server"
)

func main() {
	vaultConfig := &config.VaultConfig{
		Address:    os.Getenv("VAULT_ADDR"),
		Token:      os.Getenv("VAULT_TOKEN"),
		MountPath:  os.Getenv("VAULT_MOUNT_PATH"),
		SecretPath: os.Getenv("VAULT_SECRET_PATH"),
	}

	// Load configuration from Vault
	cfg, err := config.LoadFromVault(vaultConfig)
	if err != nil {
		log.Fatalf("Failed to load configuration from Vault: %v", err)
	}

	// Initialize dependency injection container
	container, err := container.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}
	defer container.Close()

	// Create and configure HTTP server
	srv := server.NewServer(container)

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %d", cfg.Server.Port)
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
