package repository

import (
	"github.com/nikolaevv/my-investor/internal/models"
	"gorm.io/gorm"
)

type UserDB struct {
	DB *gorm.DB
}

func NewUserDB(db *gorm.DB) *UserDB {
	return &UserDB{DB: db}
}

func (r *UserDB) Create(user *models.User) (uint, error) {
	result := r.DB.Create(user)
	return user.ID, result.Error
}

func (r *UserDB) UpdateRefreshToken(userId uint, refreshToken string) error {
	result := r.DB.Model(&models.User{}).Where("id = ?", userId).Update("refresh_token", refreshToken)
	return result.Error
}

func (r *UserDB) GetUserByLogin(login string) (*models.User, error) {
	var user models.User
	result := r.DB.Model(&models.User{}).First(&user, "login = ?", login)
	return &user, result.Error
}
