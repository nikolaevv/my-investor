package httpServer

import (
	"fmt"

	"github.com/nikolaevv/my-investor/internal/servicecontainer"
)

func New(configpath string) (*Server, error) {
	cont, err := servicecontainer.New(configpath)
	if err != nil {
		return nil, err
	}

	return &Server{
		cont: cont,
	}, nil
}

type Server struct {
	cont *servicecontainer.Container
}

func (s *Server) Start() error {
	s.configureRouter()
	s.cont.Logger.Info("Starting API server")

	return s.cont.Router.Run(
		fmt.Sprintf("%s:%s", s.cont.Config.App.Host, s.cont.Config.App.Port),
	)
}

func (s *Server) configureRouter() {
	s.cont.Router.GET("/share", s.cont.Handler.GetShare)
}
