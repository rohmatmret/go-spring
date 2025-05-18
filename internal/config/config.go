package config

import "os"

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Cache    CacheConfig
	Vault    VaultConfig
}

type DatabaseConfig struct {
	URL string
}

type ServerConfig struct {
	Host string
	Port int
}

type CacheConfig struct {
	Type  string
	Redis RedisConfig
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

func LoadConfig() (*Config, error) {
	vaultConfig := &VaultConfig{
		Address:    getEnvOrDefault("VAULT_ADDR", "http://localhost:8200"),
		Token:      getEnvOrDefault("VAULT_TOKEN", ""),
		MountPath:  getEnvOrDefault("VAULT_MOUNT_PATH", "secret"),
		SecretPath: getEnvOrDefault("VAULT_SECRET_PATH", "go-spring/config"),
	}

	return LoadFromVault(vaultConfig)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
