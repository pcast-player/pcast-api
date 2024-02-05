package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"pcast-api/model"
)

func New() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("pcast.db"), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	return db
}

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(&model.Feed{})

	if err != nil {
		panic("Failed to migrate database!")
	}
}
