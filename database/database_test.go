// database/database_test.go
package database

import (
	"testing"

	"github.com/lgarciac1603/favorites-api/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DatabaseTestSuite struct {
	suite.Suite
}

func (suite *DatabaseTestSuite) TestInitDB_ValidConnection() {
	cfg := config.DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		Database: "apidb",
		User:     "apiuser_test",
		Password: "apipass_test",
	}

	err := InitDB(cfg)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), DB)

	CloseDB()
}

func (suite *DatabaseTestSuite) TestInitDB_InvalidHost() {
	cfg := config.DatabaseConfig{
		Host:     "invalid-host-xyz-123",
		Port:     "5432",
		Database: "apidb",
		User:     "apiuser_test",
		Password: "apipass_test",
	}

	err := InitDB(cfg)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "error connecting to database")
}

func (suite *DatabaseTestSuite) TestInitDB_InvalidPassword() {
	cfg := config.DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		Database: "apidb",
		User:     "apiuser_test",
		Password: "wrong-password",
	}

	err := InitDB(cfg)

	assert.Error(suite.T(), err)
}

func (suite *DatabaseTestSuite) TestCloseDB_NoConnection() {
	DB = nil

	err := CloseDB()

	assert.NoError(suite.T(), err)
}

func TestDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(DatabaseTestSuite))
}
