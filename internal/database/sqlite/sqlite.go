package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"vault-bot/internal/database/sqldb"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "modernc.org/sqlite"
)

// Sqlite3 struct with *sql.DB instance.
// It has methods for working with URLs.
type Sqlite3 struct {
	sqldb.SQLStorage
}

// New Sqlite3 struct constructor.
func New(db *sql.DB, path string) (*Sqlite3, error) {
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return nil, fmt.Errorf("can't init migrate instance: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(path, "sqlite", driver)
	if err != nil {
		return nil, fmt.Errorf("can't create migrate instance: %w", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("can't migrate up: %w", err)
	}

	return &Sqlite3{SQLStorage: sqldb.SQLStorage{DB: db}}, nil
}
