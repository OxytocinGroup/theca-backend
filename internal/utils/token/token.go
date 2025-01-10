package token

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", nil
	}
	return base64.URLEncoding.EncodeToString(b), nil
}