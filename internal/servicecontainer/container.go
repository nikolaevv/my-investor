package servicecontainer

import (
	"github.com/gin-gonic/gin"
	"github.com/nikolaevv/my-investor/internal/handlers"
	"github.com/nikolaevv/my-investor/internal/repository"
	"github.com/nikolaevv/my-investor/pkg/auth"
	"github.com/nikolaevv/my-investor/pkg/config"
	"github.com/nikolaevv/my-investor/pkg/hash"
	"github.com/sirupsen/logrus"
)

func New(filename string) (*Container, error) {
	cfg, err := config.LoadConfig(filename)
	if err != nil {
		return nil, err
	}

	conn, err := repository.NewConnection(cfg)
	if err != nil {
		return nil, err
	}

	repo := repository.NewRepository(conn)

	hasher := hash.NewHasher()
	authManager := auth.NewAuth()

	hand, err := handlers.NewHandler(cfg, &handlers.Instruments{
		Repo:   repo,
		Hasher: hasher,
		Auth:   authManager,
	})
	if err != nil {
		return nil, err
	}

	return &Container{
		Config:  cfg,
		Logger:  logrus.New(),
		Handler: hand,
		Router:  gin.Default(),
		Repo:    repo,
		Hasher:  hasher,
		Auth:    authManager,
	}, nil
}

type Container struct {
	Config  *config.Config
	Logger  *logrus.Logger
	Handler *handlers.Handler
	Router  *gin.Engine
	Repo    *repository.Repository
	Hasher  *hash.Hasher
	Auth    *auth.Authentication
}
