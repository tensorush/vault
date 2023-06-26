package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"vault/internal/db/sqldb"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// Postgres is a struct with *sql.DB instance.
type Postgres struct {
	sqldb.SQLStore
}

// New Postgres struct constructor.
func New(db *sql.DB, path string) (*Postgres, error) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("can't init migrate instance: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(path, "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("can't create migrate instance: %w", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("can't migrate up: %w", err)
	}

	return &Postgres{SQLStore: sqldb.SQLStore{DB: db}}, nil
}
