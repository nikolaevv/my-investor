package handlers

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/nikolaevv/my-investor/internal/models"
	"github.com/nikolaevv/my-investor/internal/repository"
	mock_repository "github.com/nikolaevv/my-investor/internal/repository/mocks"
	"github.com/nikolaevv/my-investor/pkg/auth"
	mock_auth "github.com/nikolaevv/my-investor/pkg/auth/mocks"
	"github.com/nikolaevv/my-investor/pkg/config"
	"github.com/nikolaevv/my-investor/pkg/hash"
	mock_hash "github.com/nikolaevv/my-investor/pkg/hash/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHandler_login(t *testing.T) {
	type mockGetUserByLogin func(mockUser *mock_repository.MockUser, user *models.User)
	type mockHashPassword func(mockHash *mock_hash.MockPasswords, password string)
	type mockCreateAccessToken func(mockJWT *mock_auth.MockJWT, signingKey string)
	type mockCreateRefreshToken func(mockJWT *mock_auth.MockJWT)
	type mockSetRefreshToken func(mockUser *mock_repository.MockUser)
	type mockCheckPassword func(mockHash *mock_hash.MockPasswords, password string, passwordHash string)

	testTable := []struct {
		name                   string
		inputBody              string
		user                   models.User
		mockGetUserByLogin     mockGetUserByLogin
		mockSetRefreshToken    mockSetRefreshToken
		mockCreateAccessToken  mockCreateAccessToken
		mockCreateRefreshToken mockCreateRefreshToken
		mockHashPassword       mockHashPassword
		mockCheckPassword      mockCheckPassword
		inputPassword          string
		expectedStatusCode     int
	}{
		{
			name:          "OK",
			inputBody:     `{"login": "login", "password": "password"}`,
			inputPassword: "password",
			user: models.User{
				ID:           1,
				Login:        "login",
				PasswordHash: "$2a$10$8nVK6jUVdrHbRRlogTdnZ.td18pu31pkV.eqq3hTwSxP3J3chfpS2",
			},
			mockGetUserByLogin: func(mockUser *mock_repository.MockUser, user *models.User) {
				mockUser.EXPECT().GetUserByLogin(user.Login).Return(user, nil)
			},
			mockCreateAccessToken: func(mockAuth *mock_auth.MockJWT, signingKey string) {
				mockAuth.EXPECT().CreateAccessToken(1, auth.AccessTokenExpireDuration, signingKey).Return("access_token", nil)
			},
			mockCreateRefreshToken: func(mockAuth *mock_auth.MockJWT) {
				mockAuth.EXPECT().CreateRefreshToken().Return("refresh_token", nil)
			},
			mockSetRefreshToken: func(mockUser *mock_repository.MockUser) {
				mockUser.EXPECT().UpdateRefreshToken(uint(1), "refresh_token").Return(nil)
			},
			mockHashPassword: func(mockHash *mock_hash.MockPasswords, password string) {
				mockHash.EXPECT().HashAndSalt(password).Return("$2a$10$8nVK6jUVdrHbRRlogTdnZ.td18pu31pkV.eqq3hTwSxP3J3chfpS2")
			},
			mockCheckPassword: func(mockHash *mock_hash.MockPasswords, password string, passwordHash string) {
				mockHash.EXPECT().CheckPassword(password, passwordHash).Return(nil)
			},
			expectedStatusCode: 200,
		},
		{
			name:          "ERROR",
			inputBody:     `{"login": "login", "password": "password"}`,
			inputPassword: "password",
			user: models.User{
				ID:           1,
				Login:        "login",
				PasswordHash: "$2a$10$8nVK6jUVdrHbRRlogTdnZ.td18pu31pkV.eqq3hTwSxP3J3chfpS2",
			},
			mockGetUserByLogin: func(mockUser *mock_repository.MockUser, user *models.User) {
				mockUser.EXPECT().GetUserByLogin(user.Login).Return(nil, errors.New("record not found"))
			},
			mockCreateAccessToken: func(mockAuth *mock_auth.MockJWT, signingKey string) {
			},
			mockCreateRefreshToken: func(mockAuth *mock_auth.MockJWT) {
			},
			mockSetRefreshToken: func(mockUser *mock_repository.MockUser) {
			},
			mockHashPassword: func(mockHash *mock_hash.MockPasswords, password string) {
			},
			mockCheckPassword: func(mockHash *mock_hash.MockPasswords, password string, passwordHash string) {
			},
			expectedStatusCode: 404,
		},
		{
			name:          "ERROR",
			inputBody:     `{"login": "login", "password": "password"}`,
			inputPassword: "password",
			user: models.User{
				ID:           1,
				Login:        "login",
				PasswordHash: "$2a$10$8nVK6jUVdrHbRRlogTdnZ.td18pu31pkV.eqq3hTwSxP3J3chfpS2",
			},
			mockGetUserByLogin: func(mockUser *mock_repository.MockUser, user *models.User) {
				mockUser.EXPECT().GetUserByLogin(user.Login).Return(user, nil)
			},
			mockCreateAccessToken: func(mockAuth *mock_auth.MockJWT, signingKey string) {
			},
			mockCreateRefreshToken: func(mockAuth *mock_auth.MockJWT) {
			},
			mockSetRefreshToken: func(mockUser *mock_repository.MockUser) {
			},
			mockHashPassword: func(mockHash *mock_hash.MockPasswords, password string) {
			},
			mockCheckPassword: func(mockHash *mock_hash.MockPasswords, password string, passwordHash string) {
				mockHash.EXPECT().CheckPassword(password, passwordHash).Return(errors.New("password is incorrect"))
			},
			expectedStatusCode: 403,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_repository.NewMockUser(c)
			testCase.mockGetUserByLogin(user, &testCase.user)

			cfg, err := config.LoadConfig(*ConfigPath)
			if err != nil {
				panic(err)
			}

			passwordsHasher := mock_hash.NewMockPasswords(c)
			testCase.mockCheckPassword(passwordsHasher, testCase.inputPassword, testCase.user.PasswordHash)

			JWTAuth := mock_auth.NewMockJWT(c)
			testCase.mockCreateAccessToken(JWTAuth, cfg.Auth.JWTSecret)
			testCase.mockCreateRefreshToken(JWTAuth)

			testCase.mockSetRefreshToken(user)

			hasher := &hash.Hasher{Passwords: passwordsHasher}
			repository := &repository.Repository{User: user}
			authManager := &auth.Authentication{JWT: JWTAuth}

			handler, _ := NewHandler(cfg, &Instruments{
				Repo:   repository,
				Auth:   authManager,
				Hasher: hasher,
			})

			url := "/login"
			r := gin.Default()
			r.POST(url, handler.Login)

			// Test request
			httpTestRecorder := httptest.NewRecorder()
			testRequest := httptest.NewRequest("POST", url,
				bytes.NewBufferString(testCase.inputBody),
			)

			// Perform Request
			r.ServeHTTP(httpTestRecorder, testRequest)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, httpTestRecorder.Code)
		})
	}

}
