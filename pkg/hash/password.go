package hash

import (
	"golang.org/x/crypto/bcrypt"
)

type PasswordsHasher struct {
}

func NewPasswordsHasher() *PasswordsHasher {
	return &PasswordsHasher{}
}

func (h *PasswordsHasher) HashAndSalt(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

func (h *PasswordsHasher) CheckPassword(password string, passwordHash string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
}
