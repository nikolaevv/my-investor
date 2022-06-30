package usecase

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/nikolaevv/my-investor/internal/domain/entity"
	serviceСontainer "github.com/nikolaevv/my-investor/internal/domain/service/container"
	"github.com/nikolaevv/my-investor/internal/domain/service/repository"
	mock_repository "github.com/nikolaevv/my-investor/internal/domain/service/repository/mocks"
	"github.com/nikolaevv/my-investor/pkg/auth"
	mock_auth "github.com/nikolaevv/my-investor/pkg/auth/mocks"
	"github.com/nikolaevv/my-investor/pkg/config"
	mock_hash "github.com/nikolaevv/my-investor/pkg/hash/mocks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestHandler_login(t *testing.T) {
	type mockGetUserByLogin func(mockUser *mock_repository.MockUser, user *entity.User)
	type mockHashPassword func(mockHash *mock_hash.MockPasswords, password string)
	type mockCreateAccessToken func(mockJWT *mock_auth.MockJWT, signingKey string)
	type mockCreateRefreshToken func(mockJWT *mock_auth.MockJWT)
	type mockSetRefreshToken func(mockUser *mock_repository.MockUser)
	type mockCheckPassword func(mockHash *mock_hash.MockPasswords, password string, passwordHash string)

	testTable := []struct {
		name                   string
		inputBody              string
		user                   entity.User
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
			user: entity.User{
				ID:           1,
				Login:        "login",
				PasswordHash: "$2a$10$8nVK6jUVdrHbRRlogTdnZ.td18pu31pkV.eqq3hTwSxP3J3chfpS2",
			},
			mockGetUserByLogin: func(mockUser *mock_repository.MockUser, user *entity.User) {
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
			user: entity.User{
				ID:           1,
				Login:        "login",
				PasswordHash: "$2a$10$8nVK6jUVdrHbRRlogTdnZ.td18pu31pkV.eqq3hTwSxP3J3chfpS2",
			},
			mockGetUserByLogin: func(mockUser *mock_repository.MockUser, user *entity.User) {
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
			expectedStatusCode: 403,
		},
		{
			name:          "ERROR",
			inputBody:     `{"login": "login", "password": "password"}`,
			inputPassword: "password",
			user: entity.User{
				ID:           1,
				Login:        "login",
				PasswordHash: "$2a$10$8nVK6jUVdrHbRRlogTdnZ.td18pu31pkV.eqq3hTwSxP3J3chfpS2",
			},
			mockGetUserByLogin: func(mockUser *mock_repository.MockUser, user *entity.User) {
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

			cfg, err := config.LoadConfig(RelativeConfigPath)
			if err != nil {
				panic(err)
			}

			// Mock data
			user := mock_repository.NewMockUser(c)
			testCase.mockGetUserByLogin(user, &testCase.user)
			passwordsHasher := mock_hash.NewMockPasswords(c)
			testCase.mockCheckPassword(passwordsHasher, testCase.inputPassword, testCase.user.PasswordHash)
			JWTAuth := mock_auth.NewMockJWT(c)
			testCase.mockCreateAccessToken(JWTAuth, cfg.GetString("Auth.JWTSecret"))
			testCase.mockCreateRefreshToken(JWTAuth)
			testCase.mockSetRefreshToken(user)
			repository := &repository.Repository{User: user}

			r := gin.Default()
			container := &serviceСontainer.Container{
				Config: cfg,
				Logger: logrus.New(),
				Router: r,
				Repo:   repository,
				Hasher: passwordsHasher,
				Auth:   JWTAuth,
			}
			handler := NewHandler(container)

			url := "/login"
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
