package database

import (
	"database/sql"
	"fmt"
	"go_finance/pkg/utils"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
)

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
func ConnectAndMigrateDB(dataSourceName, migrationPath string) (*sql.DB, *migrate.Migrate, error) {
	db, err := ConnectDB(dataSourceName)
	if err != nil {
		return nil, nil, err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, nil, fmt.Errorf("could not create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationPath,
		"postgres",
		driver,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("migration init error: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, nil, fmt.Errorf("migration up error: %w", err)
	}

	fmt.Println("Database connected and migrations applied")
	return db, m, nil
}

func DownDB(m *migrate.Migrate) {
	if err := m.Down(); err != nil {
		utils.Logger.Error("migration failed to make down ", "error", err)
	}
}
