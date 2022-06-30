package serviceСontainer

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nikolaevv/my-investor/internal/domain/entity"
	"github.com/nikolaevv/my-investor/internal/domain/service/repository"
	"github.com/nikolaevv/my-investor/pkg/auth"
	"github.com/nikolaevv/my-investor/pkg/config"
	"github.com/nikolaevv/my-investor/pkg/database"
	"github.com/nikolaevv/my-investor/pkg/hash"
	"github.com/sirupsen/logrus"
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

func New(filename string) (*Container, error) {
	cfg, err := config.LoadConfig(filename)
	if err != nil {
		return nil, err
	}

	conn, err := database.NewConnection(cfg,
		&entity.User{},
		&entity.Share{},
	)

	if err != nil {
		return nil, err
	}

	repo := repository.NewRepository(conn)
	hasher := hash.NewPasswordsHasher()
	authManager := auth.NewJWTManager()

	return &Container{
		Config: cfg,
		Logger: logrus.New(),
		Router: gin.Default(),
		Repo:   repo,
		Hasher: hasher,
		Auth:   authManager,
	}, nil
}

type Container struct {
	Config *config.Config
	Logger *logrus.Logger
	Router *gin.Engine
	Repo   Repository
	Hasher PasswordsHasher
	Auth   JWTManager
}
