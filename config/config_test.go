// config/config_test.go
package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
	originalEnv map[string]string
}

func (suite *ConfigTestSuite) SetupTest() {
	suite.originalEnv = make(map[string]string)
	envVars := []string{"DB_HOST", "DB_PORT", "DB_NAME", "DB_USER", "DB_PASS", "APP_PORT", "AUTH_API_URL"}

	for _, key := range envVars {
		suite.originalEnv[key] = os.Getenv(key)
	}
}

func (suite *ConfigTestSuite) TeardownTest() {
	for key, val := range suite.originalEnv {
		if val == "" {
			os.Unsetenv(key)
		} else {
			os.Setenv(key, val)
		}
	}
}

func (suite *ConfigTestSuite) TestLoadConfig_WithEnvVars() {
	os.Setenv("DB_HOST", "custom-host")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_NAME", "custom-db")
	os.Setenv("DB_USER", "custom-user")
	os.Setenv("DB_PASS", "custom-pass")
	os.Setenv("APP_PORT", "9090")
	os.Setenv("AUTH_API_URL", "http://auth.local")

	cfg := LoadConfig()

	assert.Equal(suite.T(), "custom-host", cfg.Host)
	assert.Equal(suite.T(), "5433", cfg.Port)
	assert.Equal(suite.T(), "custom-db", cfg.Database)
	assert.Equal(suite.T(), "custom-user", cfg.User)
	assert.Equal(suite.T(), "custom-pass", cfg.Password)
	assert.Equal(suite.T(), "9090", cfg.AppPort)
	assert.Equal(suite.T(), "http://auth.local", cfg.AuthAPI)
}

func (suite *ConfigTestSuite) TestLoadConfig_WithDefaults() {
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASS")
	os.Unsetenv("APP_PORT")
	os.Unsetenv("AUTH_API_URL")

	cfg := LoadConfig()

	assert.Equal(suite.T(), "localhost", cfg.Host)
	assert.Equal(suite.T(), "5432", cfg.Port)
	assert.Equal(suite.T(), "apidb", cfg.Database)
	assert.Equal(suite.T(), "apiuser_test", cfg.User)
	assert.Equal(suite.T(), "apipass_test", cfg.Password)
	assert.Equal(suite.T(), "8090", cfg.AppPort)
	assert.Equal(suite.T(), "http://localhost:8080", cfg.AuthAPI)
}

func (suite *ConfigTestSuite) TestGetConnectionString() {
	cfg := DatabaseConfig{
		Host:     "localhost",
		Port:     "8090",
		Database: "testdb",
		User:     "testuser",
		Password: "testpass",
	}

	connStr := cfg.GetConnectionString()

	expected := "host=localhost port=8090 user=testuser password=testpass dbname=testdb sslmode=disable"
	assert.Equal(suite.T(), expected, connStr)
}

func (suite *ConfigTestSuite) TestGetConnectionString_WithSpecialChars() {
	cfg := DatabaseConfig{
		Host:     "db.example.com",
		Port:     "5433",
		Database: "mydb",
		User:     "admin",
		Password: "p@ss!word#123",
	}

	connStr := cfg.GetConnectionString()

	assert.Contains(suite.T(), connStr, "host=db.example.com")
	assert.Contains(suite.T(), connStr, "password=p@ss!word#123")
	assert.Contains(suite.T(), connStr, "dbname=mydb")
}

func (suite *ConfigTestSuite) TestConfigStruct_AllFieldsPresent() {
	cfg := DatabaseConfig{
		Host:     "host",
		Port:     "port",
		Database: "db",
		User:     "user",
		Password: "pass",
		AppPort:  "8090",
		AuthAPI:  "http://auth.local",
	}

	assert.NotEmpty(suite.T(), cfg.Host)
	assert.NotEmpty(suite.T(), cfg.Port)
	assert.NotEmpty(suite.T(), cfg.Database)
	assert.NotEmpty(suite.T(), cfg.User)
	assert.NotEmpty(suite.T(), cfg.Password)
	assert.NotEmpty(suite.T(), cfg.AppPort)
	assert.NotEmpty(suite.T(), cfg.AuthAPI)
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
