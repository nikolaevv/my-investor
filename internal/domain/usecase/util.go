package usecase

import (
	"context"

	serviceСontainer "github.com/nikolaevv/my-investor/internal/domain/service/container"
	"github.com/nikolaevv/my-investor/pkg/auth"
	"github.com/nikolaevv/my-investor/pkg/config"
	"github.com/nikolaevv/my-investor/pkg/gen/proto/tinkoff/investapi"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

func CreateUserSession(userId uint, repo serviceСontainer.Repository, authManager *auth.Authentication, cfg *config.Config) (Tokens, error) {
	var (
		result Tokens
		err    error
	)

	signingKey := cfg.Auth.JWTSecret

	result.AccessToken, err = authManager.JWT.CreateAccessToken(int(userId), auth.AccessTokenExpireDuration, signingKey)
	if err != nil {
		return result, err
	}

	refreshToken, err := authManager.JWT.CreateRefreshToken()
	if err != nil {
		return result, err
	}

	result.RefreshToken = refreshToken
	return result, repo.UpdateRefreshToken(userId, refreshToken)
}

func CreateTinkoffSandboxAccount(URL string, Token string, ctx context.Context) (string, error) {
	sandboxClient := investapi.CreateSandboxServiceClient(URL, Token)
	openAccountReq := investapi.OpenSandboxAccountRequest{}
	protoOpenAccountMsg, err := sandboxClient.OpenSandboxAccount(ctx, &openAccountReq)
	if err != nil {
		return "", err
	}

	return protoOpenAccountMsg.AccountId, nil
}
