package usecase

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/nikolaevv/my-investor/internal/domain/entity"
	"github.com/nikolaevv/my-investor/internal/transport/dto/request"
)

func prepareSignUpReq(c *gin.Context) (*request.UserAuth, error) {
	var req request.UserAuth
	if err := c.BindJSON(&req); err != nil {
		return nil, err
	}
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return &req, nil
}

func registerUser(h *handler, reqData *request.UserAuth, ctx context.Context) (uint, error) {
	user := &entity.User{}
	err := copier.Copy(user, reqData)
	if err != nil {
		return 0, nil
	}

	// creating sandbox stock account
	accountId, err := CreateTinkoffSandboxAccount(h.Config.Tinkoff.URL, h.Config.Tinkoff.Token, ctx)
	if err != nil {
		return 0, nil
	}

	user.AccountID = accountId
	user.PasswordHash = h.Hasher.Passwords.HashAndSalt(reqData.Password)

	return h.Repo.User.CreateUser(user)
}

func (h *handler) SignUp(c *gin.Context) {
	reqData, err := prepareSignUpReq(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	userId, err := registerUser(h, reqData, c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	tokens, err := CreateUserSession(userId, h.Repo, h.Auth, h.Config)
	if err != nil {
		c.String(http.StatusBadGateway, err.Error())
		return
	}

	c.JSON(http.StatusCreated, tokens)
}
