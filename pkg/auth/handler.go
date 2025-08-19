package auth

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"service-record/pkg/auth/templates"
)

// LoginHandler - обработчик для отображения страницы входа
func LoginHandler(config Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return templates.Login().Render(c.Context(), c.Response().BodyWriter())
	}
}

// LoginPostHandler - обработчик для обработки формы входа
func LoginPostHandler(config Config) fiber.Handler {
	service := &AuthService{DB: config.DB, Logger: config.Logger}
	return func(c *fiber.Ctx) error {
		// 1. Получаем данные из формы
		username := c.FormValue("username")
		password := c.FormValue("password")

		// 2. Проверяем данные пользователя
		user, err := service.GetUserByUsername(username)
		if err != nil || !CheckPasswordHash(password, user.PasswordHash) {
			config.Logger.Info().Msg("Invalid login attempt")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid username or password"})
		}

		// 3. Создаём новую сессию
		sessionID, err := service.CreateSession(user.ID, c.IP(), c.Get("User-Agent"))
		if err != nil {
			config.Logger.Error().Err(err).Msg("Failed to create session")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
		}

		// 4. Устанавливаем куки
		c.Cookie(&fiber.Cookie{
			Name:     config.SessionCookieName,
			Value:    sessionID,
			Expires:  time.Now().Add(config.SessionDuration),
			HTTPOnly: true,
			Path:     "/",
			SameSite: "Lax",
		})

		// 5. Перенаправляем на защищённый маршрут
		return c.Redirect("/profile")
	}
}

// LogoutHandler - обработчик для выхода
func LogoutHandler(config Config) fiber.Handler {
	service := &AuthService{DB: config.DB, Logger: config.Logger}
	return func(c *fiber.Ctx) error {
		sessionID := c.Cookies(config.SessionCookieName)
		if sessionID != "" {
			if err := service.DeleteSession(sessionID); err != nil {
				config.Logger.Error().Err(err).Msg("Failed to delete session")
			}
		}

		// Удаляем куки
		c.Cookie(&fiber.Cookie{
			Name:     config.SessionCookieName,
			Value:    "",
			Expires:  time.Now().Add(-1 * time.Hour),
			HTTPOnly: true,
			Path:     "/",
		})

		return c.Redirect(config.LoginPath)
	}
}