package auth

import (
	"fmt"
	"strconv"
	"strings"
)

type userCrypto struct {
	key []byte
}

func newUserCrypto(key []byte) *userCrypto {
	return &userCrypto{
		key: key,
	}
}

func (u *userCrypto) encrypt(userCookie *userCookie) (string, error) {
	return fmt.Sprintf("%s*%v*%v", userCookie.ID, userCookie.IsActive, userCookie.IsAdmin), nil
}

func (u *userCrypto) decrypt(encrypted string) (*userCookie, error) {
	if encrypted == "" {
		return nil, fmt.Errorf("строка для расшифровки пуста")
	}
	
	encryptedSlice := strings.Split(encrypted, "*")

	if len(encryptedSlice) != 3 {
		return nil, fmt.Errorf("некорректная строка для расшифровки")
	}
	
	isActive, err := strconv.ParseBool(encryptedSlice[1])
	if err != nil {
		return nil, fmt.Errorf("некорректное значение активности пользователя")

	}

	isAdmin, err := strconv.ParseBool(encryptedSlice[2])
	if err != nil {
		return nil, fmt.Errorf("некорректное значение признака админа")
	}
	
	return &userCookie{
		ID:       encryptedSlice[0],
		IsActive: isActive,
		IsAdmin:  isAdmin,
	}, nil
}