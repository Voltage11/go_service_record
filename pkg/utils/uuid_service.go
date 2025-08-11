package utils

import (
	"github.com/google/uuid"

	
)

func NewUUID() string {	
	return uuid.New().String()
}

func IsValidUUID(uuidStr string) bool {
	if len(uuidStr) != 36 { 
        return false
    }
    _, err := uuid.Parse(uuidStr)
    return err == nil	
}