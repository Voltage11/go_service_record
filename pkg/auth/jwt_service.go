package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims - кастомная структура claims с нашими данными пользователя
type Claims struct {
    User UserCookie `json:"user"`
    jwt.RegisteredClaims
}


// GenerateToken создает JWT токен с данными пользователя
func (a *AuthService) generateToken(user UserCookie) (string, error) {
    // Создаем claims с данными пользователя
    claims := Claims{
        User: user,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Токен действителен 24 часа
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            Issuer:    "go-service",
        },
    }

    // Создаем токен с claims
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    // Подписываем токен секретным ключом
    tokenString, err := token.SignedString(a.JwtKey)
    if err != nil {
        return "", fmt.Errorf("ошибка подписи токена: %w", err)
    }

    return tokenString, nil
}

// ParseToken проверяет и парсит JWT токен
func (a *AuthService) parseToken(tokenString string) (*UserCookie, error) {
    // Парсим токен с нашими claims
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        // Проверяем метод подписи
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
        }
        return a.JwtKey, nil
    })

    if err != nil {
        return nil, fmt.Errorf("ошибка парсинга токена: %w", err)
    }

    // Извлекаем claims
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return &claims.User, nil
    }

    return nil, fmt.Errorf("невалидный токен")
}