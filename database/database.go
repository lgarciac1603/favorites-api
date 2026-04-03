package database

import (
	"database/sql"
	"fmt"

	"github.com/lgarciac1603/favorites-api/config"
	_ "github.com/lib/pq"
)

var DB *sql.DB

// DatabaseInterface define mock operations
type DatabaseInterface interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// InitDB init DB connection
func InitDB(cfg config.DatabaseConfig) error {
	connStr := cfg.GetConnectionString()
	
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error abriendo BD: %w", err)
	}
	
	err = db.Ping()
	if err != nil {
		return fmt.Errorf("error conectando a BD: %w", err)
	}
	
	DB = db
	fmt.Println("✓ Conectado a PostgreSQL exitosamente")
	return nil
}

func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
