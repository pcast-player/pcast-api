package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"pcast-api/model"
	"time"
)

func New() *gorm.DB {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
		},
	)

	db, err := gorm.Open(sqlite.Open("./pcast.db"), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic("Failed to connect to database!")
	}

	return db
}

func NewTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("./../pcast-test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(&model.Feed{})
	if err != nil {
		panic("Failed to migrate database!")
	}
}

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
