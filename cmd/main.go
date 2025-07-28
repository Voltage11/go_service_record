package main

import (
	"fmt"
	"service-record/internal/config"
	"service-record/pkg/database"
	"service-record/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	fmt.Println("Загрузка конфигурации...")
	cfg := config.GetConfig()
	fmt.Println("Конфигурация загружена")

	
	loggerCfg := logger.New(cfg.Log.Level)
	logger := loggerCfg.Logger
	logger.Info().Msg(fmt.Sprintf("Создан логгер, уровень логирования будет изменен перед запуском сервера на: %d", loggerCfg.LogLevel))

	logger.Info().Msg("Подключение к БД...")
	db, err := database.NewDatabase(cfg.Db, *logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Ошибка подключения к БД")
	}
	defer db.Close()

	app := fiber.New()
	app.Use(recover.New())
	app.Static("/static", "./static")

	serverRun := fmt.Sprintf("%s:%s", cfg.Server.ServerHost, cfg.Server.ServerPort)

	logger.Info().Msg(fmt.Sprintf("Запуск сервера на %s", serverRun))
	logger.Info().Msg(fmt.Sprintf("Уровень логирования изменен на: %d", loggerCfg.LogLevel))
	loggerCfg.SetLogLevel()
	app.Listen(serverRun)
}