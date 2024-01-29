package models

import (
	"gorm.io/driver/sqlite"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func ConnectDatabase() {
	database, err := gorm.Open(sqlite.Open("pcast.db"), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	mErr := database.AutoMigrate(&Feed{})

	if mErr != nil {
		panic("Failed to connect to database!")
	}

	db = database
}

func GetFeeds() []Feed {
	var feeds []Feed

	db.Find(&feeds)

	return feeds
}

func CreateFeed(feed *Feed) {
	db.Create(feed)
}
