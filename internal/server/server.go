package server

import (
	"fmt"
	"net/http"

	"github.com/nikolaevv/my-investor/internal/servicecontainer"
	"github.com/rs/cors"
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

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"POST", "GET", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		Debug:            true,
	})

	modifiedHandler := c.Handler(s.cont.Router)

	s.cont.Logger.Info("Starting API server")
	return http.ListenAndServe(fmt.Sprintf("%s:%s", s.cont.Config.App.Host, s.cont.Config.App.Port), modifiedHandler)
}

func (s *Server) configureRouter() {

}
