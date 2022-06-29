package repository

import (
	"github.com/nikolaevv/my-investor/internal/models"
	"gorm.io/gorm"
)

type ShareDB struct {
	DB *gorm.DB
}

func NewShareDB(db *gorm.DB) *ShareDB {
	return &ShareDB{DB: db}
}

func (r *ShareDB) Create(share *models.Share) (uint, error) {
	result := r.DB.Create(share)
	return share.ID, result.Error
}
