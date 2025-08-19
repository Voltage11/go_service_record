package auth

import (
	"crypto/rand"
	"encoding/hex"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword - хеширование пароля
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash - проверка пароля
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateSessionID - генерация уникального ID для сессии
func GenerateSessionID() string {
	b := make([]byte, 25)
	if _, err := rand.Read(b); err != nil {
		log.Println("Error generating session ID:", err)
		return ""
	}
	return hex.EncodeToString(b)
}
