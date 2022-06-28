package handlers

import (
	"github.com/nikolaevv/my-investor/pkg/config"
)

type Handler struct {
	//Repo       *repository.Repository
	Config *config.Config
}

func New(config *config.Config) (*Handler, error) {
	/*
		repo, err := repository.New(config)
		if err != nil {
			return nil, err
		}
	*/

	return &Handler{
		//Repo:       repo,
		Config: config,
	}, nil
}
