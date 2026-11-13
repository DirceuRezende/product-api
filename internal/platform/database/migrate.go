package database

import (
	"errors"
	"fmt"
	"strings"

	"products-api/db"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// MigrateUp aplica todas as migrations pendentes embutidas no binário.
func MigrateUp(databaseURL string) error {
	src, err := iofs.New(db.Migrations, "migrations")
	if err != nil {
		return fmt.Errorf("load migrations: %w", err)
	}

	// o driver do golang-migrate para pgx v5 usa o scheme "pgx5"
	url := databaseURL
	for _, prefix := range []string{"postgres://", "postgresql://"} {
		if strings.HasPrefix(url, prefix) {
			url = "pgx5://" + strings.TrimPrefix(url, prefix)
			break
		}
	}

	m, err := migrate.NewWithSourceInstance("iofs", src, url)
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}
