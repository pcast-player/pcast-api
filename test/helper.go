package test

import (
	"gorm.io/gorm"
	"log"
	"pcast-api/domain/feed/model"
)

func TruncateTables(db *gorm.DB) {
	err := db.Exec("DELETE FROM feeds;").Error
	if err != nil {
		log.Fatal(err)
	}
}

func RemoveTables(db *gorm.DB) {
	err := db.Migrator().DropTable(&model.Feed{})
	if err != nil {
		log.Fatal(err)
	}
}
