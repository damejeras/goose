package db

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/**.sql
var migrations embed.FS

func Open(ctx context.Context, logger *slog.Logger, path string) (c *sql.DB, rerr error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("create sqlite connection: %w", err)
	}
	defer func() {
		if rerr != nil {
			db.Close()
		}
	}()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return nil, fmt.Errorf("create migration source: %w", err)
	}

	dest, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return nil, fmt.Errorf("create migration destination: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", source, "sqlite3", dest)
	if err != nil {
		return nil, fmt.Errorf("create migration: %w", err)
	}

	logger.Debug("Running database migrations", "path", path)

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("run migration: %w", err)
	} else if errors.Is(err, migrate.ErrNoChange) {
		logger.Info("No database migrations to apply")
	} else {
		logger.Info("Database migrations applied successfully")
	}

	return db, nil
}
