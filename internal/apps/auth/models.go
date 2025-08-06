package auth

import "time"

type User struct {
	ID       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	HashPass string `json:"hash_pass" db:"hash_pass"`
	IsAdmin  bool   `json:"is_admin" db:"is_admin"`
	IsActive bool   `json:"is_active" db:"is_active"`
	CreateAt string `json:"create_at" db:"create_at"`
	UpdateAt string `json:"update_at" db:"update_at"`
}

func (u *User) IsValid() bool {
	return u.Email != "" && u.HashPass != "" && u.IsActive
}

type Session struct {
	ID        string    `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	IP        string    `json:"ip" db:"ip"`
	UserAgent string    `json:"user_agent" db:"user_agent"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
func (s *Session) IsValid() bool {
	if s.ID == "" || s.UserID == 0 || s.ExpiresAt.IsZero() || s.IP == "" || s.UserAgent == "" {
		return false
	}
	return true	
}