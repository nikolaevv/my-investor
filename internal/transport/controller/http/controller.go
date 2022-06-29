package httpController

import (
	serviceСontainer "github.com/nikolaevv/my-investor/internal/domain/service/container"
	"github.com/nikolaevv/my-investor/internal/domain/usecase"
)

func ConfigureRouter(container *serviceСontainer.Container) {
	handler := usecase.NewHandler(container)
	container.Router.GET("/share", handler.GetShare)
	container.Router.POST("/signup", handler.SignUp)
	container.Router.POST("/login", handler.Login)
	container.Router.POST("/orderShare", handler.BuyShare)
}
