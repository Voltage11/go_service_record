package auth

import "time"

type User struct {
	Id           int    `json:"id" db:"id"`
	UserName     string `json:"username" db:"username"`
	PasswordHash string `json:"password_hash" db:"password_hash"`
	Email        string `json:"email" db:"email"`
	IsActive     bool   `json:"is_active" db:"is_active"`
	IsAdmin      bool   `json:"is_admin" db:"is_admin"`
}

type UserCookie struct {
	Id           int    `json:"id" db:"id"`
	UserName     string `json:"username" db:"username"`	
	Email        string `json:"email" db:"email"`
	IsActive     bool   `json:"is_active" db:"is_active"`
	IsAdmin      bool   `json:"is_admin" db:"is_admin"`
}

// type UserCookie struct {
// 	Id           int    `json:"id" db:"id"`
// 	UserName     string `json:"username" db:"username"`
// 	PasswordHash string `json:"password_hash" db:"password_hash"`
// 	Email        string `json:"email" db:"email"`
// 	IsActive     bool   `json:"is_active" db:"is_active"`
// 	IsAdmin      bool   `json:"is_admin" db:"is_admin"`
// }


type session struct {
	Id        string    `json:"id" db:"id"`
	UserId    int       `json:"user_id" db:"user_id"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	Ip        string    `json:"ip" db:"ip"`
	UserAgent string    `json:"user_agent" db:"user_agent"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
