package repository

import (
	"github.com/nikolaevv/my-investor/internal/domain/entity"
	"gorm.io/gorm"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go
type User interface {
	CreateUser(user *entity.User) (uint, error)
	UpdateRefreshToken(userId uint, refreshToken string) error
	GetUserByLogin(login string) (*entity.User, error)
	GetUserByID(id int) (*entity.User, error)
}

type Share interface {
	CreateShare(share *entity.Share) (uint, error)
}

type Repository struct {
	User
	Share
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User:  NewUser(db),
		Share: NewShare(db),
	}
}
