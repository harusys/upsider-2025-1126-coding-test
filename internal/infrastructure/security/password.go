package security

import (
	"github.com/harusys/super-shiharai-kun/internal/infrastructure"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a password using bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), infrastructure.BcryptCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// CheckPassword checks if the password matches the hash.
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
