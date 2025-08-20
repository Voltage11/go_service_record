package auth

import (
	"github.com/google/uuid"

	
)

func newUUID() string {	
	return uuid.New().String()
}

func isValidUUID(uuidStr string) bool {
	if len(uuidStr) != 36 { 
        return false
    }
    _, err := uuid.Parse(uuidStr)
    return err == nil	
}