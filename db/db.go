package db

import (
	"fmt"

	"github.com/gayathriad/go_api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=password dbname=goapi port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	db.AutoMigrate(&models.User{})
	fmt.Println("connected to database!")
	return db, nil
}
