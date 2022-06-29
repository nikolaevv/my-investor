package hash

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=password.go -destination=mocks/mock.go
type Passwords interface {
	HashAndSalt(password string) string
	CheckPassword(password string, passwordHash string) error
}

type Hasher struct {
	Passwords
}

func NewHasher() *Hasher {
	return &Hasher{
		Passwords: NewPasswordsHasher(),
	}
}

type PasswordsHasher struct {
}

func NewPasswordsHasher() *PasswordsHasher {
	return &PasswordsHasher{}
}

func (h *PasswordsHasher) HashAndSalt(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}
	return string(hash)
}

func (h *PasswordsHasher) CheckPassword(password string, passwordHash string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
}
