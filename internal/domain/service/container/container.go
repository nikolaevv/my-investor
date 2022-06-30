package serviceСontainer

import (
	"github.com/gin-gonic/gin"
	"github.com/nikolaevv/my-investor/internal/domain/entity"
	"github.com/nikolaevv/my-investor/internal/domain/service/repository"
	"github.com/nikolaevv/my-investor/pkg/auth"
	"github.com/nikolaevv/my-investor/pkg/config"
	"github.com/nikolaevv/my-investor/pkg/database"
	"github.com/nikolaevv/my-investor/pkg/hash"
	"github.com/sirupsen/logrus"
)

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
	hasher := hash.NewHasher()
	authManager := auth.NewAuth()

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
	Repo   *repository.Repository
	Hasher *hash.Hasher
	Auth   *auth.Authentication
}
