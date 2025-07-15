package main

import (
	"os"
	"service-record/internal/config"
	"service-record/internal/repository"

	"github.com/rs/zerolog"
)

func main() {
	// Инициализация логгера
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Загрузка env окружения
	logger.Info().Msg("Загрузка конфигурации...")
	cfg := config.GetConfig(&logger)
	logger.Info().Msg("Конфигурация успешно загружена")
	
	// Подключаемся к БД
	logger.Info().Msg("Подключение к БД...")
	db := repository.GetDatabase(cfg.DatabaseURL, "migrations", &logger)
	logger.Info().Msg("Подключение к БД успешно")
	defer db.Close()
	
}
