package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type AuthHandler struct {
	router fiber.Router
	logger *zerolog.Logger
}

func NewAuthHandler(router fiber.Router, customLogger *zerolog.Logger) {
	h := &AuthHandler{
		router:       router,
		logger: customLogger,
	}

	authGroup := h.router.Group("/auth")
	authGroup.Get("/login", h.getLogin).Name("auth.login")
	authGroup.Get("/logout", h.getLogout).Name("auth.logout")

	authGroup.Get("/user/list", h.getUserList).Name("auth.user-list")
}

func (h *AuthHandler) getLogin(c *fiber.Ctx) error {
	return c.SendString("login")
}

func (h *AuthHandler) getLogout(c *fiber.Ctx) error {
	return c.SendString("LOGOUT")
}

func (h *AuthHandler) getUserList(c *fiber.Ctx) error {
	return c.SendString("Пользователи")
}

// func (h *AuthHandler) isAuth(c *fiber.Ctx) bool {
// 	// проверим на наличие куков

	
// }