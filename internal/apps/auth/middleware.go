package auth

import (
	"service-record/pkg/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AuthService struct {
	jwtService *JWTService
}

func NewAuthService(jwtSecret string) *AuthService {
	return &AuthService{
		jwtService: NewJWTService(jwtSecret),
	}
}

func (a *AuthService) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionID := c.Cookies("sessionID")
		if sessionID == "" || !utils.IsValidUUID(sessionID) {
			sessionID = a.handleSessionCreation(c)
		}

		c.Locals("sessionID", sessionID)
		go a.logSessionActivity(sessionID, c.IP(), c.Path())

		if a.isPublicPath(c.Path()) {
			return c.Next()
		}

		return a.handleProtectedRoutes(c)
	}
}

func (a *AuthService) CreateJWTCookie(user *User) (*fiber.Cookie, error) {
	token, err := a.jwtService.GenerateToken(user.ID, user.IsAdmin, user.IsActive, 24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &fiber.Cookie{
		Name:     "jwtToken",
		Value:    token,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
	}, nil
}

func (a *AuthService) isPublicPath(path string) bool {
	publicPaths := map[string]bool{
		"/auth/login": true,
	}
	return publicPaths[path]
}

func (a *AuthService) setCookie(c *fiber.Ctx, key, value string) {
	c.Cookie(&fiber.Cookie{
		Name:     key,
		Value:    value,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
	})
}

func (a *AuthService) handleSessionCreation(c *fiber.Ctx) string {
	newSessionID := utils.NewUUID()
	a.setCookie(c, "sessionID", newSessionID)
	return newSessionID
}

func (a *AuthService) logSessionActivity(sessionID, ip, path string) {
	// Логирование активности
}

func (a *AuthService) handleProtectedRoutes(c *fiber.Ctx) error {
	jwtToken := c.Cookies("jwtToken")
	if jwtToken == "" {
		return c.Redirect("/auth/login", fiber.StatusTemporaryRedirect)
	}

	claims, err := a.jwtService.ParseToken(jwtToken)
	if err != nil {
		c.ClearCookie("jwtToken")
		return c.Redirect("/auth/login", fiber.StatusTemporaryRedirect)
	}

	c.Locals("user", &User{
		ID:       claims.UserID,
		IsAdmin:  claims.IsAdmin,
		IsActive: claims.IsActive,
	})

	return c.Next()
}