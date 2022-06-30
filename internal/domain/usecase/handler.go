package usecase

import (
	serviceСontainer "github.com/nikolaevv/my-investor/internal/domain/service/container"
	"github.com/nikolaevv/my-investor/pkg/config"
)

const (
	RelativeConfigPath = "../../../configs/app.json"
)

type handler struct {
	Repo   serviceСontainer.Repository
	Config config.Config
	Hasher serviceСontainer.PasswordsHasher
	Auth   serviceСontainer.JWTManager
	Logger serviceСontainer.Logger
}

func NewHandler(container *serviceСontainer.Container) *handler {
	return &handler{
		Config: container.Config,
		Repo:   container.Repo,
		Hasher: container.Hasher,
		Auth:   container.Auth,
		Logger: container.Logger,
	}
}
