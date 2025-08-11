package auth

import (
	"context"
	"errors"

	"service-record/pkg/utils"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

const (
	cookieSessionName = "sessionID"
	cookieCryptoToken = "ctoken"
)

type AuthService struct {
	userCrypto *userCrypto
	logger     *zerolog.Logger
	cryptoKey  []byte
	repository *repository
	hashKey    string
}

func NewAuthService(db *sqlx.DB, logger *zerolog.Logger, cryptoKey []byte, hashKey string) *AuthService {
	return &AuthService{
		userCrypto: newUserCrypto(cryptoKey),
		logger:     logger,
		cryptoKey:  cryptoKey,
		repository: newRepository(db, logger),
		hashKey:    hashKey,
	}
}

func (a *AuthService) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Get("User-Agent") == "" {
			return c.Status(403).SendString("User-Agent header не указан!")
		}

		sessionID := c.Cookies(cookieSessionName)
		if sessionID == "" || !utils.IsValidUUID(sessionID) {
			sessionID = a.handleSessionCreation(c)
		}

		c.Locals(cookieSessionName, sessionID)
		go a.logSessionActivity(sessionID, c.IP(), c.Path())

		if a.isPublicPath(c.Path()) {
			return c.Next()
		}

		return a.handleProtectedRoutes(c, sessionID)
	}
}

// Каждому пользователю создаем
func (a *AuthService) createCryptoCookie(user *User) (*fiber.Cookie, error) {
	token, err := a.userCrypto.encrypt(nil)
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

// Публичные пути, которые не требуют авторизации
func (a *AuthService) isPublicPath(path string) bool {
	// Если url начинается на /public, то пропускаем всегда
	partsPath := strings.Split(path, "/")
	if len(partsPath) > 0 {
		if partsPath[0] == "public" {
			return true
		}
	}

	publicPaths := map[string]bool{
		"/auth/login":     true,
		"/auth/api/login": true,
	}

	return publicPaths[path]
}

// Установка куки
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

// Кажому пользователю генерируется уникальный токен сессии, если отсутствует
func (a *AuthService) handleSessionCreation(c *fiber.Ctx) string {
	newSessionID := utils.NewUUID()
	a.setCookie(c, cookieSessionName, newSessionID)

	return newSessionID
}

func (a *AuthService) logSessionActivity(sessionID, ip, path string) {
	// Логирование активности
}

func (a *AuthService) handleProtectedRoutes(c *fiber.Ctx, sessionID string) error {
	cryptoToken := c.Cookies(cookieCryptoToken)
	if cryptoToken == "" {
		return c.Redirect("/auth/login", fiber.StatusTemporaryRedirect)
	}

	claims, err := a.userCrypto.decrypt(cryptoToken)
	if err != nil {
		c.ClearCookie(cookieCryptoToken)
		return c.Redirect("/auth/login", fiber.StatusTemporaryRedirect)
	}

	if !a.isAuthValid(claims.ID, sessionID) {

	}

	c.Locals("user", &User{
		ID:       claims.ID,
		IsAdmin:  claims.IsAdmin,
		IsActive: claims.IsActive,
	})

	return c.Next()
}

func (a *AuthService) isAuthValid(userId string, sessionID string) bool {
	if userId == "" || sessionID == "" {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	userSessions, err := a.repository.getUserSessions(ctx, userId)
	if err != nil || len(userSessions) == 0 {
		return false
	}

	for _, session := range userSessions {
		if session.ID == sessionID {
			return true
		}
	}
	return false
}

// Авторизация пользователя, если все верно, то сосздадим куки
func (a *AuthService) authUser(c *fiber.Ctx, email, password, sessionID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// fmt.Println(utils.StrToHashWithKey("admin", a.hashKey))
	// fmt.Println(utils.NewUUID())

	// Получим пользователя по email
	user, err := a.repository.getByEmail(ctx, email)
	if err != nil {
		return errors.New("Пользователь не найден по email: " + email)
	}
	if !user.IsActive {
		return errors.New("Пользователь не активный: " + email)
	}
	// Зашифруем пароль
	passwordHash := utils.StrToHashWithKey(password, a.hashKey)

	if user.PasswordHash != passwordHash {
		return errors.New("Неверный пароль")
	}

	// Если все верно, то создаем куки
	cryptoUser, err := a.userCrypto.encrypt(&userCookie{
		ID:       user.ID,
		IsAdmin:  user.IsAdmin,
		IsActive: user.IsActive,
	})
	if err != nil {
		a.logger.Error().Err(err).Msg("Ошибка шифрования: " + email)
		return errors.New("Ошибка шифрования: " + email)
	}
	// Добавим в базу сессию пользователя
	session := &session{
		ID:       sessionID,
		UserID:   user.ID,
		IP:       c.IP(),
		UserAgent: string(c.Request().Header.UserAgent()),
		CreatedAt: time.Now(),
	}
	err = a.repository.createSession(ctx, session)
	if err != nil {
		a.logger.Error().Err(err).Msg("Ошибка создания сессии в бд")
		return errors.New("Ошибка создания сессии в бд: " + sessionID)
	}

	a.setCookie(c, cookieCryptoToken, cryptoUser)

	return nil

}
