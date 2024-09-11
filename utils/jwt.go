package utils

import (
	"crypto/rand"
	"encoding/base64"
	"os"
)

var JWTSecret []byte

func InitJWTSecret() error {
	secretStr := os.Getenv("JWT_SECRET")
	if secretStr == "" {
		// Generate a random secret if not provided
		secret := make([]byte, 32)
		_, err := rand.Read(secret)
		if err != nil {
			return err
		}
		secretStr = base64.StdEncoding.EncodeToString(secret)
		os.Setenv("JWT_SECRET", secretStr)
	}
	JWTSecret = []byte(secretStr)
	return nil
}
