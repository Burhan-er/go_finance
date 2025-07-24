package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // Postgres driver
)

// ConnectDB establishes a connection to the database
func ConnectDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("could not open sql connection: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("could not ping database: %w", err)
	}

	return db, nil
}