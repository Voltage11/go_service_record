package main

import (
	"fmt"
	"os"
	"service-record/internal/config"

	"github.com/rs/zerolog"
)

func main() {
	// Инициализация логгера
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Загрузка env окружения
	logger.Info().Msg("Загрузка конфигурации...")
	cfg := config.GetConfig(&logger)
	logger.Info().Msg("Конфигурация успешно загружена")
	fmt.Println(cfg)
}
