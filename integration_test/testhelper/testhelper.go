package testhelper

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"

	"github.com/labstack/echo/v4"
	"pcast-api/config"
	"pcast-api/controller"
	"pcast-api/db"
	"pcast-api/router"
)

var DB *sql.DB

const TestDSN = "host=localhost port=5432 user=pcast password=pcast dbname=pcast_test sslmode=disable"

func Setup() {
	DB = db.NewTestDB(TestDSN)

	DB.Exec("TRUNCATE TABLE users CASCADE")
	DB.Exec("TRUNCATE TABLE feeds")
	DB.Exec("TRUNCATE TABLE episodes")

	RunMigrations()
}

func Teardown() {
	DB.Exec("TRUNCATE TABLE users CASCADE")
	DB.Exec("TRUNCATE TABLE feeds")
	DB.Exec("TRUNCATE TABLE episodes")
	DB.Close()
}

func RunMigrations() {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS episodes (
			id UUID PRIMARY KEY,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			feed_id UUID NOT NULL,
			feed_guid VARCHAR(255) NOT NULL,
			current_position INTEGER,
			played BOOLEAN NOT NULL DEFAULT FALSE
		)
	`)
	if err != nil {
		log.Printf("Warning: episodes table creation: %v\n", err)
	}
	DB.Exec(`CREATE INDEX IF NOT EXISTS idx_episodes_feed_id ON episodes(feed_id)`)
	DB.Exec(`CREATE INDEX IF NOT EXISTS idx_episodes_feed_guid ON episodes(feed_guid)`)

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS feeds (
			id UUID PRIMARY KEY,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			user_id UUID NOT NULL,
			title VARCHAR(500) NOT NULL,
			url VARCHAR(1000) NOT NULL,
			synced_at TIMESTAMP
		)
	`)
	if err != nil {
		log.Printf("Warning: feeds table creation: %v\n", err)
	}
	DB.Exec(`CREATE INDEX IF NOT EXISTS idx_feeds_user_id ON feeds(user_id)`)

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL
		)
	`)
	if err != nil {
		log.Panicf("CRITICAL: failed to create users table: %v", err)
	}
	DB.Exec(`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`)
}

func NewApp() *echo.Echo {
	r := router.NewTestRouter()
	apiGroup := r.Group("/api")

	cfg := &config.Config{
		Auth: config.Auth{
			JwtSecret:        "testsecret",
			JwtExpirationMin: 10,
		},
	}

	controller.NewController(cfg, DB, apiGroup)

	return r
}

func Unmarshal[T any](bytes []byte) (*T, error) {
	m := new(T)
	err := json.Unmarshal(bytes, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func UnmarshalResult[T any](result io.Reader) (*T, error) {
	bytes, err := io.ReadAll(result)
	if err != nil {
		return nil, err
	}
	return Unmarshal[T](bytes)
}

func TruncateAll() {
	DB.Exec("TRUNCATE TABLE users CASCADE")
	DB.Exec("TRUNCATE TABLE feeds")
	DB.Exec("TRUNCATE TABLE episodes")
}
