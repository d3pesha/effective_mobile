package database

import (
	"em/internal/config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormHandler(cfg *config.Config) (*gorm.DB, error) {
	ds := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DatabaseHost,
		cfg.DatabaseUser,
		cfg.DatabasePassword,
		cfg.DatabaseName,
		cfg.DatabasePort,
	)

	db, err := gorm.Open(postgres.Open(ds), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
