package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"pcast-api/config"
)

// New returns a standard sql.DB connection for use with sqlc
func New(c *config.Config) *sql.DB {
	dsn := c.Database.GetPostgresDSN()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %v", err))
	}

	if err := createConnectionPool(c, db); err != nil {
		panic(err)
	}

	return db
}

func createConnectionPool(c *config.Config, db *sql.DB) error {
	duration, err := time.ParseDuration(c.Database.MaxLifetime)
	if err != nil {
		return err
	}

	db.SetMaxIdleConns(c.Database.MaxIdleConnections)
	db.SetMaxOpenConns(c.Database.MaxConnections)
	db.SetConnMaxLifetime(duration)

	return nil
}

// NewTestDB returns a sql.DB connection for tests using Postgres
func NewTestDB(dsn string) *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(fmt.Sprintf("failed to connect test database: %v", err))
	}

	return db
}
