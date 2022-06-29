package app

import (
	"fmt"
	serviceСontainer "github.com/nikolaevv/my-investor/internal/domain/service/container"
	httpController "github.com/nikolaevv/my-investor/internal/transport/controller/http"
)

func New(configpath string) (*Server, error) {
	cont, err := serviceСontainer.New(configpath)
	if err != nil {
		return nil, err
	}

	return &Server{
		cont: cont,
	}, nil
}

type Server struct {
	cont *serviceСontainer.Container
}

func (s *Server) Start() error {
	httpController.ConfigureRouter(s.cont)
	s.cont.Logger.Info("Starting API server")

	return s.cont.Router.Run(
		fmt.Sprintf("%s:%s", s.cont.Config.App.Host, s.cont.Config.App.Port),
	)
}
