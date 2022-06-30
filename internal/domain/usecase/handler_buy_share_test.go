package usecase

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/nikolaevv/my-investor/internal/domain/entity"
	serviceСontainer "github.com/nikolaevv/my-investor/internal/domain/service/container"
	"github.com/nikolaevv/my-investor/internal/domain/service/repository"
	mock_repository "github.com/nikolaevv/my-investor/internal/domain/service/repository/mocks"
	"github.com/nikolaevv/my-investor/pkg/auth"
	mock_auth "github.com/nikolaevv/my-investor/pkg/auth/mocks"
	"github.com/nikolaevv/my-investor/pkg/config"
	"github.com/nikolaevv/my-investor/pkg/hash"
	mock_hash "github.com/nikolaevv/my-investor/pkg/hash/mocks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestHandler_buyShare(t *testing.T) {
	type mockAuthorizateUser func(mockJWT *mock_auth.MockJWT, headers http.Header, claims *auth.Claims, signingKey string)
	type mockGetUserById func(mockUser *mock_repository.MockUser, user *entity.User)
	type mockCreateShare func(mockShare *mock_repository.MockShare, share *entity.Share)

	testTable := []struct {
		name                string
		inputBody           string
		user                *entity.User
		share               *entity.Share
		mockAuthorizateUser mockAuthorizateUser
		mockGetUserById     mockGetUserById
		mockCreateShare     mockCreateShare
		expectedStatusCode  int
	}{
		{
			name:      "OK",
			inputBody: `{"id": "GAZP", "classCode": "TQBR", "quantity": 1}`,
			user: &entity.User{
				ID: 1,
			},
			share: &entity.Share{
				Code:      "GAZP",
				ClassCode: "TQBR",
				UserID:    1,
				Quantity:  1,
			},
			mockAuthorizateUser: func(mockJWT *mock_auth.MockJWT, headers http.Header, claims *auth.Claims, signingKey string) {
				mockJWT.EXPECT().AuthorizateUser(headers, signingKey).Return(claims, nil)
			},
			mockCreateShare: func(mockShare *mock_repository.MockShare, share *entity.Share) {
				mockShare.EXPECT().CreateShare(share).Return(uint(1), nil)
			},
			mockGetUserById: func(mockUser *mock_repository.MockUser, user *entity.User) {
				mockUser.EXPECT().GetUserByID(int(user.ID)).Return(user, nil)
			},
			expectedStatusCode: 200,
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

			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			accountId, err := CreateTinkoffSandboxAccount(cfg.Tinkoff.URL, cfg.Tinkoff.Token, ctx)
			if err != nil {
				panic(err)
			}
			testCase.user.AccountID = accountId

			claims := &auth.Claims{
				Id: int(testCase.user.ID),
			}
			mockHeaders := http.Header{}

			mockJWTAuth := mock_auth.NewMockJWT(c)
			testCase.mockAuthorizateUser(mockJWTAuth, mockHeaders, claims, cfg.Auth.JWTSecret)

			mockShare := mock_repository.NewMockShare(c)
			testCase.mockCreateShare(mockShare, testCase.share)

			mockUser := mock_repository.NewMockUser(c)
			testCase.mockGetUserById(mockUser, testCase.user)

			passwordsHasher := mock_hash.NewMockPasswords(c)

			hasher := &hash.Hasher{Passwords: passwordsHasher}
			repository := &repository.Repository{User: mockUser, Share: mockShare}
			authManager := &auth.Authentication{JWT: mockJWTAuth}

			URL := "/share/order"
			r := gin.Default()

			container := &serviceСontainer.Container{
				Config: cfg,
				Logger: logrus.New(),
				Router: r,
				Repo:   repository,
				Hasher: hasher,
				Auth:   authManager,
			}
			handler := NewHandler(container)
			r.POST(URL, handler.BuyShare)

			// Test request
			httpTestRecorder := httptest.NewRecorder()
			testRequest := httptest.NewRequest("POST", URL,
				bytes.NewBufferString(testCase.inputBody),
			)

			// Perform Request
			r.ServeHTTP(httpTestRecorder, testRequest)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, httpTestRecorder.Code)
		})
	}
}
