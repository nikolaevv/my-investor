package entity

import "gorm.io/gorm"

type Share struct {
	gorm.Model
	Code      string
	ClassCode string
	User      User
	UserID    uint
	Quantity  int
}
