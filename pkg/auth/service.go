package auth

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

// AuthService - структура для бизнес-логики аутентификации
type AuthService struct {
	DB     *sqlx.DB
	Logger *zerolog.Logger
}

// GetUserByUsername - поиск пользователя по имени
func (s *AuthService) GetUserByUsername(username string) (*User, error) {
	user := &User{}
	err := s.DB.Get(user, "SELECT * FROM users WHERE username = $1", username)
	if err != nil {
		s.Logger.Error().Err(err).Msg("Failed to get user by username")
		return nil, err
	}
	return user, nil
}

// GetUserByID - поиск пользователя по ID
func (s *AuthService) GetUserByID(id int) (*User, error) {
	user := &User{}
	err := s.DB.Get(user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		s.Logger.Error().Err(err).Msg("Failed to get user by id")
		return nil, err
	}
	return user, nil
}

// GetUserByAPIKey - поиск пользователя по API ключу (для примера, так как в таблице нет такого поля, используем username)
func (s *AuthService) GetUserByAPIKey(apiKey string) (*User, error) {
	user := &User{}
	err := s.DB.Get(user, "SELECT * FROM users WHERE is_admin = true AND username = $1", apiKey)
	if err != nil {
		s.Logger.Error().Err(err).Msg("Failed to get admin user by api key")
		return nil, err
	}
	return user, nil
}

// CreateSession - создание новой сессии в базе данных
func (s *AuthService) CreateSession(userID int, ip, userAgent string) (string, error) {
	sessionID := GenerateSessionID()
	expiresAt := time.Now().Add(time.Hour * 24 * 7) // 7 дней
	_, err := s.DB.Exec("INSERT INTO sessions (id, user_id, expires_at, ip, user_agent) VALUES ($1, $2, $3, $4, $5)",
		sessionID, userID, expiresAt, ip, userAgent)
	if err != nil {
		s.Logger.Error().Err(err).Msg("Failed to create session")
		return "", err
	}
	return sessionID, nil
}

// GetSession - получение сессии по ID
func (s *AuthService) GetSession(sessionID string) (*Session, error) {
	session := &Session{}
	err := s.DB.Get(session, "SELECT * FROM sessions WHERE id = $1", sessionID)
	if err != nil {
		s.Logger.Error().Err(err).Msg("Failed to get session")
		return nil, err
	}
	return session, nil
}

// DeleteSession - удаление сессии по ID
func (s *AuthService) DeleteSession(sessionID string) error {
	_, err := s.DB.Exec("DELETE FROM sessions WHERE id = $1", sessionID)
	if err != nil {
		s.Logger.Error().Err(err).Msg("Failed to delete session")
		return err
	}
	return nil
}
