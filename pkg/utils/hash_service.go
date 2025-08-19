package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

)

func StrToHashWithKey(strToHash string, key string) string {	
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(strToHash))
	
	//fmt.Println(hex.EncodeToString(h.Sum(nil)), strToHash, key)
	return hex.EncodeToString(h.Sum(nil))
}