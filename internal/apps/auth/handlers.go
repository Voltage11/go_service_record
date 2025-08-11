package auth

import (
	"net/http"
	"service-record/pkg/tadapter"	
	"service-record/views/components"
	"service-record/views/pages"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type AuthHandler struct {
	router      fiber.Router
	logger      *zerolog.Logger
	authService *AuthService
}

func NewAuthHandler(router fiber.Router, customLogger *zerolog.Logger, authService *AuthService) {
	h := &AuthHandler{
		router:      router,
		logger:      customLogger,
		authService: authService,
	}

	authGroup := h.router.Group("/auth")
	authGroup.Get("/login", h.getLogin).Name("auth.login")
	authGroup.Get("/logout", h.getLogout).Name("auth.logout")

	authGroup.Post("/api/login", h.postLogin).Name("auth.api-login")

	authGroup.Get("/user/list", h.getUserList).Name("auth.user-list")
}

func (h *AuthHandler) getLogin(c *fiber.Ctx) error {
	component := pages.Auth()

	return tadapter.Render(c, component, http.StatusOK)
}

func (h *AuthHandler) postLogin(c *fiber.Ctx) error {

	login := c.FormValue("login")
	password := c.FormValue("password")	
	sessionID := c.Cookies(cookieSessionName)
	
	err := h.authService.authUser(c, login, password, sessionID)
	if err == nil {
		return c.Redirect("/", http.StatusFound)
	}

	component := components.Message(
		components.MessageProps{
			Type:    components.MessageTypeError,
			Message: "Не верный логин/пароль",
		},
	)

	return tadapter.Render(c, component, http.StatusOK)
}

func (h *AuthHandler) getLogout(c *fiber.Ctx) error {
	return c.SendString("LOGOUT")
}

func (h *AuthHandler) getUserList(c *fiber.Ctx) error {
	return c.SendString("Пользователи")
}
