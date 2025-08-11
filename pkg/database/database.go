package database

import (
	"fmt"
	"service-record/internal/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)


func NewDatabase(cfg config.DbConfig, logger zerolog.Logger) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", cfg.Dsn)
	if err != nil {
		logger.Error().Err(err).Msg("Ошибка подключения к БД")
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	if err := db.Ping(); err != nil {
		logger.Error().Err(err).Msg("Ошибка пинга к БД")
		return nil, fmt.Errorf("ошибка пинга к БД: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	logger.Info().Msg("Запуск миграций...")
	if err := runMigrations(db, cfg.MigrationsPath, logger); err != nil {
		return nil, fmt.Errorf("ошибка миграций: %w", err)
	}

	return db, nil
}

func runMigrations(db *sqlx.DB, migrationsPath string, logger zerolog.Logger) error {
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
		return fmt.Errorf("ошибка миграции: %w", err)
	}

	// defer func() {
	// 	if _, err := m.Close(); err != nil {
	// 		logger.Error().Err(err).Msg("Ошибка закрытия миграции")
	// 	}
	// }()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("ошибка применения миграции: %w", err)
	}

	logger.Info().Msg("Миграции выполнены успешно")
	return nil
}