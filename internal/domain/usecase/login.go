package usecase

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nikolaevv/my-investor/internal/transport/dto/request"
)

func (h *handler) Login(c *gin.Context) {
	var req request.UserAuth
	if err := c.BindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("error: %s", err))
		return
	}
	if err := req.Validate(); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("error: %s", err))
		return
	}

	user, err := h.Repo.User.GetUserByLogin(req.Login)
	if err != nil {
		c.String(http.StatusNotFound, fmt.Sprintf("error: %s", err))
		return
	}

	if err = h.Hasher.CheckPassword(req.Password, user.PasswordHash); err != nil {
		c.String(http.StatusForbidden, fmt.Sprintf("error: %s", err))
		return
	}

	tokens, err := CreateUserSession(
		user.ID,
		h.Repo,
		h.Auth,
		h.Config,
	)
	if err != nil {
		c.String(http.StatusBadGateway, fmt.Sprintf("error: %s", err))
		return
	}

	c.JSON(200, tokens)
}
