package postgresql

import (
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(db *PostgresDB, migrationsPath string, logger *slog.Logger) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		logger.Error("failed to create migration driver", "error", err)
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		logger.Error("failed to create migrate instance", "error", err)
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			logger.Info("no migrations to run, database already up to date")
			return nil
		}
		logger.Error("failed to run migrations", "error", err)
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	logger.Info("database migrations applied successfully")
	return nil
}
