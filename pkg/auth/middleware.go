package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

const (
	cookieSessionName string = "sessionId"
	cookieJwtName string = "jwtToken"
	sessionTimeout time.Duration = 24 * time.Hour	

	pathAuthLogin string = "/auth/login"
	pathAuthLogout string = "/auth/logout"	
)

type AuthService struct {
	Db *sqlx.DB
	Logger *zerolog.Logger
	HashKey []byte
	JwtKey []byte
	CookieSessionName string
	router fiber.Router
	repository *repository
}

func NewMiddleware(router fiber.Router, db *sqlx.DB, logger *zerolog.Logger, hashKey []byte, jwtKey []byte) (*AuthService, error) {
	if logger == nil {
		fmt.Errorf("не указан логгер при создании сервиса авторизации")
		return nil, errors.New("не указан логгер при создании сервиса авторизации")
	}
	
	if db == nil {				
		logger.Error().Msg("не указана бд при создании сервиса авторизации")		
		return nil, errors.New("не указана бд при создании сервиса авторизации")
	}
	if hashKey == nil {
		logger.Error().Msg("не указан ключ шифрования при создании сервиса авторизации")
		return nil, errors.New("не указан ключ шифрования при создании сервиса авторизации")
	}
	if jwtKey == nil {
		logger.Error().Msg("не указан ключ шифрования при создании сервиса авторизации")
		return nil, errors.New("не указан ключ шифрования при создании сервиса авторизации")		
	}
	return &AuthService{
		Db: db,
		Logger: logger,
		HashKey: hashKey,
		JwtKey: jwtKey,
		CookieSessionName: cookieSessionName,
		router: router,
		repository: newRepository(db, logger),
	}, nil
	
}

func (a *AuthService) Middleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        requestType := getRequestType(c)
		// По api пока закроем доступ, позже отработаем
		if requestType == requestTypeApi || requestType == requestTypeUnknown {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"err": "Доступ запрещен",
				"code": fiber.StatusForbidden,
			})						
		}
							
		sessionId, newSession := a.checkCookieSession(c)
		c.Locals(cookieSessionName, sessionId)
		c.Locals("isAuth", false)
		// Если новая сессия, то перенаправим на страницу авторизации сразу
		if newSession && c.Path() != pathAuthLogin {			
			return c.Redirect(pathAuthLogin)			
		}
		
		// Если метод POST на странице авторизации, то пропустим
		if c.Path() == "/auth/api/login" && c.Method() == "POST" {			
			return c.Next()
		}
		
		// Проверка JWT
		jwtToken := c.Cookies(cookieJwtName)
		if jwtToken == "" {
			if !a.isPublicPath(c) {
				return c.Redirect(pathAuthLogin)
			} else {				
				return c.Next()				
			}			
		}

		userCookie, err := a.parseToken(jwtToken)
		if err != nil {
			c.ClearCookie(cookieJwtName)
			if !a.isPublicPath(c) {
				return c.Redirect(pathAuthLogin)
			} else {				
				return c.Next()				
			}
		}
		if userCookie != nil {
			c.Locals("isAuth", true)
			c.Locals("user", userCookie)
			// Если пользователь авторизован, то перенаправим на главную страницу
			if c.Path() == pathAuthLogin {
				return c.Redirect("/", fiber.StatusFound)
			}			
		} else {
			if !a.isPublicPath(c) {
				return c.Redirect(pathAuthLogin)
			} else {				
				return c.Next()				
			}
		}
		            
        return c.Next()
    }
}

func (a *AuthService) checkCookieSession(c *fiber.Ctx) (string, bool) {
	sessionId := c.Cookies(cookieSessionName)
	newSession := false
	if sessionId == "" {
		a.setCookieNewSession(c, newUUID())
		newSession = true
	} else {
		if !isValidUUID(sessionId) {
			c.ClearCookie(cookieSessionName)
			a.setCookieNewSession(c, newUUID())
			newSession = true
		} else {
			a.setCookieNewSession(c, sessionId)
		}
	}
	return sessionId, newSession
}

func (a *AuthService) setCookieNewSession(c *fiber.Ctx, sessionId string) {	
	c.Cookie(&fiber.Cookie{
		Name:     cookieSessionName,
		Value:    sessionId,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Expires:  time.Now().Add(sessionTimeout),
		Path:     "/",
	})
}

func (a *AuthService) isPublicPath(c *fiber.Ctx) bool {
	publicPath := map[string]bool{
		pathAuthLogin: true,
		"/": true,
	}
	
	if _, ok := publicPath[c.Path()]; ok {
		return true
	}
	return false
}