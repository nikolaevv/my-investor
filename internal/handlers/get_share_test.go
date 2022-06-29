package handlers

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/nikolaevv/my-investor/internal/repository"
	mock_repository "github.com/nikolaevv/my-investor/internal/repository/mocks"
	"github.com/nikolaevv/my-investor/pkg/auth"
	mock_auth "github.com/nikolaevv/my-investor/pkg/auth/mocks"
	"github.com/nikolaevv/my-investor/pkg/config"
	"github.com/nikolaevv/my-investor/pkg/hash"
	mock_hash "github.com/nikolaevv/my-investor/pkg/hash/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHandler_getShare(t *testing.T) {
	testTable := []struct {
		name               string
		tickerId           string
		classCode          string
		expectedStatusCode int
	}{
		{
			name:               "OK",
			tickerId:           "GAZP",
			classCode:          "TQBR",
			expectedStatusCode: 200,
		},
		{
			name:               "OK",
			tickerId:           "AAAAAAAAA",
			classCode:          "TQBR",
			expectedStatusCode: 404,
		},
		{
			name:               "OK",
			tickerId:           "GAZP",
			classCode:          "BBBBBBBBB",
			expectedStatusCode: 400,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			cfg, err := config.LoadConfig(*ConfigPath)
			if err != nil {
				panic(err)
			}

			user := mock_repository.NewMockUser(c)
			passwordsHasher := mock_hash.NewMockPasswords(c)
			JWTAuth := mock_auth.NewMockJWT(c)

			hasher := &hash.Hasher{Passwords: passwordsHasher}
			repository := &repository.Repository{User: user}
			authManager := &auth.Authentication{JWT: JWTAuth}

			handler, _ := NewHandler(cfg, &Instruments{
				Repo:   repository,
				Auth:   authManager,
				Hasher: hasher,
			})

			url := "/share"
			r := gin.Default()
			r.GET(url, handler.GetShare)

			// Test request
			httpTestRecorder := httptest.NewRecorder()
			testRequest := httptest.NewRequest("GET", url, nil)

			q := testRequest.URL.Query()
			q.Add("id", testCase.tickerId)
			q.Add("classCode", testCase.classCode)
			testRequest.URL.RawQuery = q.Encode()

			// Perform Request
			r.ServeHTTP(httpTestRecorder, testRequest)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, httpTestRecorder.Code)
		})
	}
}