package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/nikolaevv/my-investor/internal/handlers/requests"
	"github.com/nikolaevv/my-investor/internal/models"
)

func (h *Handler) SignUp(c *gin.Context) {
	var req requests.UserAuth
	if err := c.BindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("error: %s", err))
		return
	}
	if err := req.Validate(); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("error: %s", err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user := &models.User{}
	copier.Copy(user, req)

	// creating sandbox stock account
	accountId, err := CreateTinkoffSandboxAccount(h.Config.Tinkoff.URL, h.Config.Tinkoff.Token, ctx)
	if err != nil {
		c.String(http.StatusBadGateway, fmt.Sprintf("error: %s", err))
		return
	}

	user.AccountID = accountId
	user.PasswordHash = h.Hasher.Passwords.HashAndSalt(req.Password)

	userId, err := h.Repo.User.Create(user)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("error: %s", err))
		return
	}

	tokens, err := CreateUserSession(
		userId,
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
