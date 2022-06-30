package repository

import (
	"github.com/nikolaevv/my-investor/internal/domain/entity"
	"gorm.io/gorm"
)

type share struct {
	db *gorm.DB
}

func NewShare(db *gorm.DB) *share {
	return &share{db}
}

func (r *share) CreateShare(share *entity.Share) (uint, error) {
	result := r.db.Create(share)
	return share.ID, result.Error
}
