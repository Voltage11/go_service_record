package database

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

type Database struct {
	db     *sqlx.DB
	logger *zerolog.Logger
}

func GetDatabase(dsn string, migrationsPath string, logger *zerolog.Logger) *Database {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		logger.Fatal().Err(err).Msg("Ошибка соединения с БД")
	}

	if err = db.Ping(); err != nil {
		logger.Fatal().Err(err).Msg("Ошибка пинга к БД")
	}

	if err := runMigrations(db, migrationsPath, logger); err != nil {
		logger.Fatal().Err(err).Msg("Ошибка миграций")
	}

	return &Database{
		db:     db,
		logger: logger,
	}
}

func (d *Database) Close() error {
	return d.db.Close()
}

func runMigrations(db *sqlx.DB, migrationsPath string, logger *zerolog.Logger) error {
	logger.Info().Msg("Запуск миграций...")

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("ошибка создания драйвера миграции: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("ошибка инициализации миграции: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("ошибка применения миграций: %w", err)
	}

	logger.Info().Msg("Миграции успешно выполнены")
	return nil
}
