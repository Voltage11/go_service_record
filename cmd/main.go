package main

import (
	"fmt"
	"time"

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

	// // Midleware
	// authService, err := auth.NewAuthService(db, logger, []byte(cfg.Secret.JwtToken), cfg.Secret.HashKey)
	// if err != nil {
	// 	logger.Fatal().Err(err).Msg("Ошибка получения сервиса авторизации")
	// }
	// app.Use(authService.Middleware())

	// // Хандлеры
	// auth.NewAuthHandler(app, logger, authService)
	// app.Get("/", func(c *fiber.Ctx) error {
		
	// 	if c.Locals("user") != nil {
	// 		return c.JSON(fiber.Map{"message": "Hello, " + c.Locals("user").(*auth.User).Email})
	// 	} else {
	// 		return c.JSON(fiber.Map{"message": "Hello, stranger", "now": time.Now()})
	// 	}
	// })

	///////////////////////
	// Настройки для пакета аутентификации
	authConfig := auth.Config{
		DB:                db,
		Logger:            logger,
		SessionSecretKey:  "SESSION_SECRET_KEY",//os.Getenv("SESSION_SECRET_KEY"),
		SessionCookieName: "sessionId",//os.Getenv("SESSION_COOKIE_NAME"),
		LoginPath:         "/auth/login",
		LoginPostPath:     "/auth/login-post",
		LogoutPath:        "/auth/logout",
		AdminAPIKeyHeader: "ADMIN_API_KEY_HEADER" ,//os.Getenv("ADMIN_API_KEY_HEADER"),
		SessionDuration:   time.Hour * 24 * 7, // 7 дней
		PublicPaths:       []string{"/auth/login", "/auth/login-post", "/public"},
	}

	// Применяем наше middleware для аутентификации
	app.Use(auth.NewAuthMiddleware(authConfig))

	// Группа маршрутов для аутентификации
	authGroup := app.Group("/auth")
	authGroup.Get("/login", auth.LoginHandler(authConfig))
	authGroup.Post("/login-post", auth.LoginPostHandler(authConfig))
	authGroup.Get("/logout", auth.LogoutHandler(authConfig))

	//////////////////////////

	serverRun := fmt.Sprintf("%s:%s", cfg.Server.ServerHost, cfg.Server.ServerPort)

	logger.Info().Msg(fmt.Sprintf("Запуск сервера на %s", serverRun))
	logger.Info().Msg(fmt.Sprintf("Уровень логирования изменен на: %d", loggerCfg.LogLevel))
	loggerCfg.SetLogLevel()
	app.Listen(serverRun)
}