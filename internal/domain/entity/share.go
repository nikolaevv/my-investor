package entity

import "gorm.io/gorm"

type Share struct {
	gorm.Model
	Ticker    string
	ClassCode string
	User      User
	UserID    uint
	Quantity  int
}
