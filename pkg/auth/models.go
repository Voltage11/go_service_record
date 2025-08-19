package auth

import (
	"database/sql"
	"time"
)

// User - структура для таблицы users
type User struct {
	ID           int       `db:"id"`
	Username     string    `db:"username"`
	PasswordHash string    `db:"password_hash"`
	Email        string    `db:"email"`
	IsActive     bool      `db:"is_active"`
	IsAdmin      bool      `db:"is_admin"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// Session - структура для таблицы sessions
type Session struct {
	ID        string         `db:"id"`
	UserID    int            `db:"user_id"`
	ExpiresAt time.Time      `db:"expires_at"`
	IP        sql.NullString `db:"ip"`
	UserAgent string         `db:"user_agent"`
	CreatedAt time.Time      `db:"created_at"`
}