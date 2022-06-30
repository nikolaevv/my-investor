package app

import (
	"context"
	"fmt"
	serviceСontainer "github.com/nikolaevv/my-investor/internal/domain/service/container"
	httpController "github.com/nikolaevv/my-investor/internal/transport/controller/http"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", s.cont.Config.App.Host, s.cont.Config.App.Port),
		Handler: s.cont.Router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.cont.Logger.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	s.cont.Logger.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		s.cont.Logger.Fatalf("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		s.cont.Logger.Info("timeout of 5 seconds.")
	}
	s.cont.Logger.Info("Server exiting")
	return nil
	//return s.cont.Router.Run()
}
