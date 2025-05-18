package config

import (
	"fmt"

	vault "github.com/hashicorp/vault/api"
)

type VaultConfig struct {
	Address    string
	Token      string
	MountPath  string
	SecretPath string
}

func LoadFromVault(vaultConfig *VaultConfig) (*Config, error) {
	// Create Vault client
	config := vault.DefaultConfig()
	config.Address = vaultConfig.Address

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	// Set token
	client.SetToken(vaultConfig.Token)

	// Read secret
	secret, err := client.Logical().Read(fmt.Sprintf("%s/data/%s", vaultConfig.MountPath, vaultConfig.SecretPath))
	if err != nil {
		return nil, fmt.Errorf("failed to read secret: %w", err)
	}

	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no data found at path: %s", vaultConfig.SecretPath)
	}

	// Extract data from secret
	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid secret data format")
	}

	// Convert to Config struct
	cfg := &Config{}

	// Database config
	if db, ok := data["database"].(map[string]interface{}); ok {
		cfg.Database.URL = db["url"].(string)
	}

	// Server config
	if server, ok := data["server"].(map[string]interface{}); ok {
		cfg.Server.Port = int(server["port"].(float64))
	}

	// Cache config
	if cache, ok := data["cache"].(map[string]interface{}); ok {
		cfg.Cache.Type = cache["type"].(string)
		if redis, ok := cache["redis"].(map[string]interface{}); ok {
			cfg.Cache.Redis.Host = redis["host"].(string)
			cfg.Cache.Redis.Port = int(redis["port"].(float64))
			cfg.Cache.Redis.Password = redis["password"].(string)
			cfg.Cache.Redis.DB = int(redis["db"].(float64))
		}
	}

	return cfg, nil
}
