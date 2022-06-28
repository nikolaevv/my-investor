package repository

import (
	"fmt"

	"github.com/nikolaevv/my-investor/internal/models"
	"github.com/nikolaevv/my-investor/pkg/config"
	"github.com/nikolaevv/my-investor/pkg/env"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func New(config *config.Config) (*Repository, error) {
	conn, err := connect(config)
	if err != nil {
		return nil, err
	}

	if err = migrate(conn); err != nil {
		return nil, err
	}

	return &Repository{
		DB: conn,
	}, nil
}

func connect(config *config.Config) (*gorm.DB, error) {
	host := env.GetHost(config.DB.Host, "db")

	dsn := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=%s password=%s",
		host,
		config.DB.User,
		config.DB.Name,
		config.DB.SSLMode,
		config.DB.Pass,
	)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func migrate(conn *gorm.DB) error {
	return conn.AutoMigrate(
		&models.Share{},
		&models.User{},
	)
}
