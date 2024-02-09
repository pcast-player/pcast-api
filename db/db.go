package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"pcast-api/config"
	"time"
)

func New(c *config.Config) *gorm.DB {
	l := getLogger(c)

	db, err := gorm.Open(sqlite.Open("./pcast.db"), &gorm.Config{Logger: l})
	if err != nil {
		panic("Failed to connect to database!")
	}

	return db
}

func getLogger(c *config.Config) logger.Interface {
	if c.Database.Logging {
		return logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: true,
				ParameterizedQueries:      true,
			},
		)
	} else {
		return nil
	}
}

func NewTestDB(dsn string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}
