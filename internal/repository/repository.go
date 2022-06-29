package repository

import (
	"github.com/nikolaevv/my-investor/internal/models"
	"gorm.io/gorm"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go
type User interface {
	Create(user *models.User) (uint, error)
	UpdateRefreshToken(userId uint, refreshToken string) error
	GetUserByLogin(login string) (*models.User, error)
}

type Repository struct {
	User
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User: NewUserDB(db),
	}
}
