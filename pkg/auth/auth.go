package auth

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"

)

// UserInfoKey - ключ для хранения данных пользователя в контексте Fiber.
type UserInfoKey struct{}

// Config - структура для настроек middleware
type Config struct {
	DB                *sqlx.DB
	Logger            *zerolog.Logger
	SessionSecretKey  string
	SessionCookieName string
	LoginPath         string
	LoginPostPath     string
	LogoutPath        string
	AdminAPIKeyHeader string
	SessionDuration   time.Duration
	PublicPaths       []string
}

// NewAuthMiddleware - основная функция, которая возвращает Fiber middleware
func NewAuthMiddleware(config Config) fiber.Handler {
	service := &AuthService{DB: config.DB, Logger: config.Logger}
	return func(c *fiber.Ctx) error {
		// Проверка, является ли маршрут публичным (не требует аутентификации)
		for _, path := range config.PublicPaths {
			if strings.HasPrefix(c.Path(), path) {
				return c.Next()
			}
		}

		// Попытка аутентификации через API Key для администраторов
		apiKey := c.Get(config.AdminAPIKeyHeader)
		if apiKey != "" {
			user, err := service.GetUserByAPIKey(apiKey)
			if err == nil && user != nil && user.IsAdmin {
				c.Locals(UserInfoKey{}, user)
				return c.Next()
			}
		}

		// Попытка аутентификации через сессию (куки)
		sessionID := c.Cookies(config.SessionCookieName)
		if sessionID != "" {
			session, err := service.GetSession(sessionID)
			if err == nil && session != nil && session.ExpiresAt.After(time.Now()) {
				user, err := service.GetUserByID(session.UserID)
				if err == nil && user != nil {
					c.Locals(UserInfoKey{}, user)
					return c.Next()
				}
			}
		}

		// Если аутентификация не удалась, перенаправляем или возвращаем ошибку
		if strings.Contains(c.Get("Accept"), "application/json") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		// Перенаправление на страницу входа
		return c.Redirect(config.LoginPath)
	}
}

// GetCurrentUser - универсальная функция для получения текущего пользователя из контекста
func GetCurrentUser(c *fiber.Ctx) *User {
	if user, ok := c.Locals(UserInfoKey{}).(*User); ok {
		return user
	}
	return nil
}