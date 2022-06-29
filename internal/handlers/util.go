package handlers

import (
	"github.com/nikolaevv/my-investor/internal/repository"
	"github.com/nikolaevv/my-investor/pkg/auth"
	"github.com/nikolaevv/my-investor/pkg/config"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

func CreateUserSession(userId uint, repo *repository.Repository, authManager *auth.Authentication, cfg *config.Config) (Tokens, error) {
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
	repo.User.UpdateRefreshToken(userId, refreshToken)

	return result, err
}
