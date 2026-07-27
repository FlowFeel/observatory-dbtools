// Package migrate handles execution of SQL migrations via golang-migrate.
package migrate

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Runner manages applying SQL migrations against a target DB.
type Runner struct {
	db        *sql.DB
	sourceURL string
}

// NewRunner returns a migration runner for the given migrations folder.
func NewRunner(db *sql.DB, migrationsDir string) *Runner {
	return &Runner{
		db:        db,
		sourceURL: fmt.Sprintf("file://%s", migrationsDir),
	}
}

// Up applies all pending up migrations.
func (r *Runner) Up() error {
	m, err := r.createMigrateInstance()
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate: up: %w", err)
	}
	return nil
}

// Down rolls back all applied migrations.
func (r *Runner) Down() error {
	m, err := r.createMigrateInstance()
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate: down: %w", err)
	}
	return nil
}

func (r *Runner) createMigrateInstance() (*migrate.Migrate, error) {
	driver, err := mysql.WithInstance(r.db, &mysql.Config{})
	if err != nil {
		return nil, fmt.Errorf("migrate: driver instance: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(r.sourceURL, "mysql", driver)
	if err != nil {
		return nil, fmt.Errorf("migrate: new instance: %w", err)
	}

	return m, nil
}
