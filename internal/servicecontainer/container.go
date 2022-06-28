package servicecontainer

import (
	"github.com/gin-gonic/gin"
	"github.com/nikolaevv/my-investor/pkg/config"
	"github.com/sirupsen/logrus"
)

func New(filename string) (*Container, error) {
	cfg, err := config.LoadConfig(filename)
	if err != nil {
		return nil, err
	}

	/*
		hand, err := handlers.New(cfg)
		if err != nil {
			return nil, err
		}
	*/

	return &Container{
		Config: cfg,
		Logger: logrus.New(),
		//Handler: hand,
		Router: gin.Default(),
	}, nil
}

type Container struct {
	Config *config.Config
	Logger *logrus.Logger
	//Handler *handlers.Handler
	Router *gin.Engine
}
