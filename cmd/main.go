package main

import (
	"fmt"

	"service-record/internal/config"
	"service-record/pkg/auth"
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

	// middleware сервис авторизации и аутентификации
	authService, err := auth.NewMiddleware(app, db, logger, []byte(cfg.Secret.HashKey), []byte(cfg.Secret.JwtToken))
	if err != nil {		
		logger.Fatal().Err(err).Msg("Ошибка получения сервиса авторизации")		
	}
	app.Use(authService.Middleware())
	authService.NewAuthHandler()
	
	
	////тест
	app.Get("/", func(c *fiber.Ctx) error {
		sessionId := c.Cookies(authService.CookieSessionName)
		isAuth := c.Locals("isAuth")
		user := c.Locals("user")
		return c.SendString("Главная страница " + sessionId + fmt.Sprintf(" | isAuth:%t", isAuth) + fmt.Sprintf(" | user:%v", user))
	})
	////тест

	serverRun := fmt.Sprintf("%s:%s", cfg.Server.ServerHost, cfg.Server.ServerPort)

	logger.Info().Msg(fmt.Sprintf("Запуск сервера на %s", serverRun))
	logger.Info().Msg(fmt.Sprintf("Уровень логирования изменен на: %d", loggerCfg.LogLevel))
	loggerCfg.SetLogLevel()
	app.Listen(serverRun)
}