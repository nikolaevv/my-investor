package handlers

import (
	"flag"

	"github.com/nikolaevv/my-investor/internal/repository"
	"github.com/nikolaevv/my-investor/pkg/auth"
	"github.com/nikolaevv/my-investor/pkg/config"
	"github.com/nikolaevv/my-investor/pkg/hash"
)

var (
	ConfigPath = flag.String("configPath", "../../configs/app.json", "path to config file")
)

type Instruments struct {
	Repo   *repository.Repository
	Hasher *hash.Hasher
	Auth   *auth.Authentication
}

type Handler struct {
	Repo   *repository.Repository
	Config *config.Config
	Hasher *hash.Hasher
	Auth   *auth.Authentication
}

func NewHandler(config *config.Config, instruments *Instruments) (*Handler, error) {
	return &Handler{
		Config: config,
		Repo:   instruments.Repo,
		Hasher: instruments.Hasher,
		Auth:   instruments.Auth,
	}, nil
}
