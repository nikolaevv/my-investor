package usecase

import (
	serviceСontainer "github.com/nikolaevv/my-investor/internal/domain/service/container"
	"github.com/nikolaevv/my-investor/internal/domain/service/repository"
	"github.com/nikolaevv/my-investor/pkg/auth"
	"github.com/nikolaevv/my-investor/pkg/config"
	"github.com/nikolaevv/my-investor/pkg/hash"
)

const (
	RelativeConfigPath = "../../../configs/app.json"
)

type handler struct {
	Repo   *repository.Repository
	Config *config.Config
	Hasher *hash.Hasher
	Auth   *auth.Authentication
}

func NewHandler(container *serviceСontainer.Container) *handler {
	return &handler{
		Config: container.Config,
		Repo:   container.Repo,
		Hasher: container.Hasher,
		Auth:   container.Auth,
	}
}
