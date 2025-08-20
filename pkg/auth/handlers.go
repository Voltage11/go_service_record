package auth

import (
	"fmt"
	"net/http"
	"service-record/pkg/tadapter"
	"strings"
	"time"

	//"service-record/views/components"
	"service-record/views/pages"

	"github.com/gofiber/fiber/v2"
)


func (a *AuthService) NewAuthHandler() {	
	authGroup := a.router.Group("/auth")
	authGroup.Get("/login", a.getLogin).Name("auth.login")
	authGroup.Get("/logout", a.getLogout).Name("auth.logout")
	authGroup.Post("/api/login", a.postLogin).Name("auth.api-login")	
}

func (a *AuthService) getLogin(c *fiber.Ctx) error {		
	message := c.Query("message")
	
	component := pages.Auth(message)
	a.Logger.Info().Msg(message)

	return tadapter.Render(c, component, http.StatusOK)
}

func (a *AuthService) postLogin(c *fiber.Ctx) error {
    login := c.FormValue("login")
    password := c.FormValue("password")
    //sessionID := c.Cookies(cookieSessionName)

	login = strings.ToLower(login)
	passwordHash := StrToHashWithKey(password, string(a.HashKey))
	
	user, err := a.repository.getUserByEmail(c.Context(), login)
	if err != nil {
		a.Logger.Error().Err(err).Msg(fmt.Sprintf("Ошибка при получении пользователя по email: %s", login))
		return c.Redirect("/auth/login?message=пользователь не найден")
	}
	if user == nil {
		return c.Redirect("/auth/login?message=пользователь не найден")
	}

	if user.Email == login && user.PasswordHash == passwordHash {
		userCookie := UserCookie{
			Id: user.Id,
			UserName: user.UserName,
			Email: user.Email,
			IsActive: user.IsActive,
			IsAdmin: user.IsAdmin,
		}
		
		jwtToken, err := a.generateToken(userCookie)
		if err != nil {
			return c.Redirect("/auth/login?message=ошибка на сервере")
		} 

		c.Cookie(&fiber.Cookie{
			Name:     cookieJwtName,
			Value:    jwtToken,
			HTTPOnly: true,
			Secure:   true,
			SameSite: "Strict",
			Expires:  time.Now().Add(sessionTimeout),
			Path:     "/",
		})
		a.Logger.Info().Msg(fmt.Sprintf("авторизовался: %s", login))
		return c.Redirect("/", http.StatusFound)
	}
	
	return c.Redirect("/auth/login?message=неверныйлогинароль")
}

func (a *AuthService) getLogout(c *fiber.Ctx) error {
	return c.SendString("LOGOUT")
}

func (a *AuthService) getUserList(c *fiber.Ctx) error {
	return c.SendString("Пользователи")
}
