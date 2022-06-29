package database

import (
	"fmt"

	"github.com/nikolaevv/my-investor/pkg/config"
	"github.com/nikolaevv/my-investor/pkg/env"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Entity interface {
}

func NewConnection(config *config.Config, entities ...interface{}) (*gorm.DB, error) {
	conn, err := connect(config)
	if err != nil {
		return nil, err
	}

	if err = migrate(conn, entities); err != nil {
		return nil, err
	}

	return conn, nil
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

func migrate(conn *gorm.DB, entities []interface{}) error {
	return conn.AutoMigrate(entities...)
}
