package security

import (
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"golang.org/x/crypto/bcrypt"
)

type bcryptHasher struct{}

func NewPasswordHasher() ports.PasswordHasher {
	return &bcryptHasher{}
}

func (h *bcryptHasher) Hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func (h *bcryptHasher) Compare(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err == nil
}
