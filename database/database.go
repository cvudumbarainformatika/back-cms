package database

import (
	"fmt"
	"time"

	"github.com/cvudumbarainformatika/backend/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Database wraps the sqlx.DB connection
type Database struct {
	DB *sqlx.DB
}

// NewDatabase creates a new database connection based on configuration
func NewDatabase(cfg config.DatabaseConfig) (*Database, error) {
	var dsn string
	var driverName string

	switch cfg.Connection {
	case "mysql":
		driverName = "mysql"
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
			cfg.Username,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.Database,
		)
	case "postgres":
		driverName = "postgres"
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.Host,
			cfg.Port,
			cfg.Username,
			cfg.Password,
			cfg.Database,
		)
	default:
		return nil, fmt.Errorf("unsupported database connection: %s", cfg.Connection)
	}

	db, err := sqlx.Connect(driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool settings from config
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	return &Database{DB: db}, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	if d.DB != nil {
		return d.DB.Close()
	}
	return nil
}

// Ping tests the database connection
func (d *Database) Ping() error {
	if d.DB != nil {
		return d.DB.Ping()
	}
	return fmt.Errorf("database connection is nil")
}
