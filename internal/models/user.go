package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Login        string `gorm:"not null,unique"`
	PasswordHash string `json:"-" gorm:"not null"`
	RefreshToken string `json:"-"`
}
