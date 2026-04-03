// config/config.go
package config

import (
	"fmt"
	"os"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
	AppPort  string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		Database: getEnv("DB_NAME", "apidb"),
		User:     getEnv("DB_USER", "apiuser_test"),
		Password: getEnv("DB_PASS", "apipass_test"),
		AppPort:  getEnv("APP_PORT", "8090"),
	}
}

// GetConnectionString generates the PostgreSQL connection string
func (c DatabaseConfig) GetConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.Database,
	)
}

// getEnv retrieves an environment variable with a default value
func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}