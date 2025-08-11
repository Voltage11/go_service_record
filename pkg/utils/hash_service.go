package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func StrToHashWithKey(strToHash string, key string) string {	
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(strToHash))
	
	
	return hex.EncodeToString(h.Sum(nil))
}