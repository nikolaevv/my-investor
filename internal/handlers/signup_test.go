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

func TestHandler_signUp(t *testing.T) {
	type mockRegister func(mockUser *mock_repository.MockUser, user *models.User)
	type mockHashPassword func(mockHash *mock_hash.MockPasswords, password string)
	type mockCreateAccessToken func(mockUser *mock_auth.MockJWT, signingKey string)
	type mockCreateRefreshToken func(mockUser *mock_auth.MockJWT)
	type mockSetRefreshToken func(mockUser *mock_repository.MockUser)

	testTable := []struct {
		name                   string
		inputBody              string
		inputUser              models.User
		mockRegister           mockRegister
		mockHashPassword       mockHashPassword
		mockSetRefreshToken    mockSetRefreshToken
		mockCreateAccessToken  mockCreateAccessToken
		mockCreateRefreshToken mockCreateRefreshToken
		inputPassword          string
		expectedStatusCode     int
	}{
		{
			name:          "OK",
			inputBody:     `{"login": "login", "password": "password"}`,
			inputPassword: "password",
			inputUser: models.User{
				Login:        "login",
				PasswordHash: "$2a$10$8nVK6jUVdrHbRRlogTdnZ.td18pu31pkV.eqq3hTwSxP3J3chfpS2",
			},
			mockRegister: func(mockUser *mock_repository.MockUser, user *models.User) {
				mockUser.EXPECT().Create(user).Return(uint(1), nil)
			},
			mockHashPassword: func(mockHash *mock_hash.MockPasswords, password string) {
				mockHash.EXPECT().HashAndSalt(password).Return("$2a$10$8nVK6jUVdrHbRRlogTdnZ.td18pu31pkV.eqq3hTwSxP3J3chfpS2")
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
			expectedStatusCode: 200,
		},
		{
			name:          "ERROR",
			inputBody:     `{"login": "login", "password": "password"}`,
			inputPassword: "password",
			inputUser: models.User{
				Login:        "login",
				PasswordHash: "$2a$10$8nVK6jUVdrHbRRlogTdnZ.td18pu31pkV.eqq3hTwSxP3J3chfpS2",
			},
			mockRegister: func(mockUser *mock_repository.MockUser, user *models.User) {
				mockUser.EXPECT().Create(user).Return(uint(0), errors.New("login is already in db"))
			},
			mockHashPassword: func(mockHash *mock_hash.MockPasswords, password string) {
				mockHash.EXPECT().HashAndSalt(password).Return("$2a$10$8nVK6jUVdrHbRRlogTdnZ.td18pu31pkV.eqq3hTwSxP3J3chfpS2")
			},
			mockCreateAccessToken: func(mockAuth *mock_auth.MockJWT, signingKey string) {

			},
			mockCreateRefreshToken: func(mockAuth *mock_auth.MockJWT) {

			},
			mockSetRefreshToken: func(mockUser *mock_repository.MockUser) {

			},
			expectedStatusCode: 400,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			passwordsHasher := mock_hash.NewMockPasswords(c)
			testCase.mockHashPassword(passwordsHasher, testCase.inputPassword)

			user := mock_repository.NewMockUser(c)
			testCase.mockRegister(user, &testCase.inputUser)

			cfg, err := config.LoadConfig(*ConfigPath)
			if err != nil {
				panic(err)
			}

			JWTAuth := mock_auth.NewMockJWT(c)
			testCase.mockCreateAccessToken(JWTAuth, cfg.Auth.JWTSecret)
			testCase.mockCreateRefreshToken(JWTAuth)

			testCase.mockSetRefreshToken(user)

			hasher := &hash.Hasher{Passwords: passwordsHasher}
			repository := &repository.Repository{User: user}
			authManager := &auth.Authentication{JWT: JWTAuth}

			handler, _ := NewHandler(cfg, &Instruments{
				Repo:   repository,
				Hasher: hasher,
				Auth:   authManager,
			})

			r := gin.Default()
			r.POST("/signup", handler.SignUp)

			// Test request
			httpTestRecorder := httptest.NewRecorder()
			testRequest := httptest.NewRequest("POST", "/signup",
				bytes.NewBufferString(testCase.inputBody),
			)

			// Perform Request
			r.ServeHTTP(httpTestRecorder, testRequest)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, httpTestRecorder.Code)
		})
	}
}
