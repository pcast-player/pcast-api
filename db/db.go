package db

import (
	"gorm.io/driver/postgres"
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
	gc := &gorm.Config{Logger: l}

	db, err := gorm.Open(postgres.Open(c.Database.GetPostgresDSN()), gc)
	if err != nil {
		panic(err)
	}

	err = createConnectionPool(c, db)
	if err != nil {
		panic(err)
	}

	return db
}

func createConnectionPool(c *config.Config, db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	duration, err := time.ParseDuration(c.Database.MaxLifetime)
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(c.Database.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(c.Database.MaxConnections)
	sqlDB.SetConnMaxLifetime(duration)

	return nil
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
