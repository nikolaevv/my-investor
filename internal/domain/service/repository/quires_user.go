package repository

import (
	"github.com/nikolaevv/my-investor/internal/domain/entity"
	"gorm.io/gorm"
)

type user struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) *user {
	return &user{db}
}

func (r *user) CreateUser(user *entity.User) (uint, error) {
	result := r.db.Create(user)
	return user.ID, result.Error
}

func (r *user) UpdateRefreshToken(userId uint, refreshToken string) error {
	result := r.db.Model(&entity.User{}).Where("id = ?", userId).Update("refresh_token", refreshToken)
	return result.Error
}

func (r *user) GetUserByLogin(login string) (*entity.User, error) {
	var user entity.User
	result := r.db.Model(&entity.User{}).First(&user, "login = ?", login)
	return &user, result.Error
}

func (r *user) GetUserByID(id int) (*entity.User, error) {
	var user entity.User
	result := r.db.Model(&entity.User{}).First(&user, "id = ?", id)
	return &user, result.Error
}
