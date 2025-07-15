package main

import (
	"fmt"
	"log"
	"service-record/internal/config"
	"service-record/pkg/database"
	"service-record/pkg/logger"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Загрузка env окружения
	log.Println("Загрузка конфигурации...")
	cfg := config.GetConfig()
	log.Println("Конфигурация успешно загружена")

	// Инициализация логгера
	appLoggerCfg := logger.NewLogger(cfg.LogLevel)
	appLoger := appLoggerCfg.Logger

	// Подключаемся к БД
	appLoger.Info().Msg("Подключение к БД...")
	db := database.GetDatabase(cfg.DatabaseURL, cfg.MigrationPath, appLoger)
	appLoger.Info().Msg("Подключение к БД успешно")
	defer db.Close()

	app := fiber.New()
	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: appLoger,
	}))
	app.Use(recover.New()) //middleware, при падении fiber не упадет приложение
	app.Static("/static", "./static")

	/* До старта приложения установим уровни логирования для отображения
	   Сразу это не сделали, чтобы видеть процесс инициализации и запуска */
	appLoggerCfg.SetLevelInfo()

	// Стартуем сервер
	appLoger.Info().Msg(fmt.Sprintf("Запуск сервера на: %s", cfg.ServerAddress))
	app.Listen(cfg.ServerAddress)

}
