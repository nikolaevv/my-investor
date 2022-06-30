package serviceСontainer

import (
	"net/http"
	"time"

	"github.com/nikolaevv/my-investor/internal/domain/service/repository"
	"github.com/nikolaevv/my-investor/pkg/auth"
)

type Repository interface {
	repository.User
	repository.Share
}

//go:generate mockgen -source=manager.go -destination=../../../../pkg/auth/mocks/mock.go
type JWTManager interface {
	CreateAccessToken(userId int, expireDuration time.Duration, signingKey string) (string, error)
	CreateRefreshToken() (string, error)
	ParseToken(accessToken string, signingKey string) (map[string]interface{}, error)
	AuthorizateUser(headers http.Header, SigningKey string) (*auth.Claims, error)
}

//go:generate mockgen -source=container.go -destination=../../../../pkg/hash/mocks/mock.go
type PasswordsHasher interface {
	HashAndSalt(password string) string
	CheckPassword(password string, passwordHash string) error
}
