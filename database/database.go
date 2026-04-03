// database/database.go
package database

import (
	"database/sql"
	"fmt"

	"github.com/lgarciac1603/favorites-api/config"
	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB initializes the connection to PostgreSQL
func InitDB(cfg config.DatabaseConfig) error {
	connStr := cfg.GetConnectionString()
	
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}
	
	// Verify that the connection is valid
	err = db.Ping()
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}
	
	DB = db
	fmt.Println("Successfully connected to PostgreSQL")
	return nil
}

// CloseDB closes the database connection
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}