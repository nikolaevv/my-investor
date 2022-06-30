package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID           uint `gorm:"primaryKey"`
	AccountID    string
	Login        string `gorm:"not null,unique"`
	PasswordHash string `json:"-" gorm:"not null"`
	RefreshToken string `json:"-"`
}
