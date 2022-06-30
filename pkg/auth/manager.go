package auth

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	AccessTokenExpireDuration = time.Hour * 10
)

type Claims struct {
	jwt.StandardClaims
	Id int `json:"id"`
}

type JWTManager struct {
}

func NewJWTManager() *JWTManager {
	return &JWTManager{}
}

func (m *JWTManager) CreateAccessToken(userId int, expireDuration time.Duration, signingKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		Id: userId,
	})

	return token.SignedString([]byte(signingKey))
}

func (m *JWTManager) CreateRefreshToken() (string, error) {
	b := make([]byte, 32)
	currentTime := rand.NewSource(time.Now().Unix())
	randomizedValue := rand.New(currentTime)

	_, err := randomizedValue.Read(b)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

func (m *JWTManager) ParseToken(accessToken string, signingKey string) (map[string]interface{}, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return []byte(signingKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("error get user claims from token")
	}

	return claims, nil
}

func (m *JWTManager) AuthorizateUser(headers http.Header, SigningKey string) (*Claims, error) {
	reqToken := headers.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) < 2 {
		return nil, fmt.Errorf("invalid Authorization header")
	}

	reqToken = splitToken[1]

	parsedData, err := m.ParseToken(reqToken, SigningKey)
	if err != nil {
		return nil, err
	}
	userId := int(parsedData["id"].(float64))

	return &Claims{Id: userId}, nil
}
