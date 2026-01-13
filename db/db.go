package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"pcast-api/config"
)

// New returns a GORM DB connection
// TODO: This will be deprecated once all stores migrate to sqlc
func New(c *config.Config) *gorm.DB {
	l := getLogger(c)
	gc := &gorm.Config{Logger: l}

	db, err := gorm.Open(postgres.Open(c.Database.GetPostgresDSN()), gc)
	if err != nil {
		panic(err)
	}

	err = createConnectionPoolGORM(c, db)
	if err != nil {
		panic(err)
	}

	return db
}

// NewSQL returns a standard sql.DB connection for use with sqlc
func NewSQL(c *config.Config) *sql.DB {
	dsn := c.Database.GetPostgresDSN()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %v", err))
	}

	if err := createConnectionPoolSQL(c, db); err != nil {
		panic(err)
	}

	return db
}

func createConnectionPoolGORM(c *config.Config, db *gorm.DB) error {
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

func createConnectionPoolSQL(c *config.Config, db *sql.DB) error {
	duration, err := time.ParseDuration(c.Database.MaxLifetime)
	if err != nil {
		return err
	}

	db.SetMaxIdleConns(c.Database.MaxIdleConnections)
	db.SetMaxOpenConns(c.Database.MaxConnections)
	db.SetConnMaxLifetime(duration)

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

// NewTestDB returns a GORM DB connection for tests
// TODO: Migrate to use Postgres instead of SQLite
func NewTestDB(dsn string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

// NewTestDBSQL returns a sql.DB connection for tests using Postgres
func NewTestDBSQL(dsn string) *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(fmt.Sprintf("failed to connect test database: %v", err))
	}

	return db
}
