package auth

import (
	"time"
)

type User struct {
	ID           string `json:"id" db:"id"`
	Name         string `json:"username" db:"username"`
	Email        string `json:"email" db:"email"`
	PasswordHash string `json:"password_hash" db:"password_hash"`
	IsAdmin      bool   `json:"is_admin" db:"is_admin"`
	IsActive     bool   `json:"is_active" db:"is_active"`
	CreateAt     string `json:"create_at" db:"create_at"`
	UpdateAt     string `json:"update_at" db:"update_at"`
}

func (u *User) IsValid() bool {
	return u.Email != "" && u.PasswordHash != "" && u.IsActive
}

// Cессии пользователей
type session struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	IP        string    `json:"ip" db:"ip"`
	UserAgent string    `json:"user_agent" db:"user_agent"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func (s *session) IsValid() bool {
	if s.ID == "" || s.UserID == "" || s.ExpiresAt.IsZero() || s.IP == "" || s.UserAgent == "" {
		return false
	}
	return true
}

// Пользователь из куков
type userCookie struct {
	ID       string `json:"id"`
	IsActive bool   `json:"is_active"`
	IsAdmin  bool   `json:"is_admin"`
}
