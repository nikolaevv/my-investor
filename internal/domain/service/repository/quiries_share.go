package repository

import (
	"github.com/nikolaevv/my-investor/internal/domain/entity"
	"gorm.io/gorm"
)

type ShareDB struct {
	db *gorm.DB
}

func NewShareDB(db *gorm.DB) *ShareDB {
	return &ShareDB{db}
}

func (r *ShareDB) CreateShare(share *entity.Share) (uint, error) {
	result := r.db.Create(share)
	return share.ID, result.Error
}
