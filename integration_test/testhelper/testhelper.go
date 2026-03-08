package testhelper

import (
	"database/sql"
	"encoding/json"
	"io"

	"github.com/labstack/echo/v4"

	"pcast-api/config"
	"pcast-api/controller"
	"pcast-api/db"
	"pcast-api/router"
)

var DB *sql.DB

// TestDSN is the connection string for the test database
const TestDSN = "host=localhost port=5432 user=pcast password=pcast dbname=pcast_test sslmode=disable"

// Test configuration constants
const (
	TestJWTSecret        = "testsecret"
	TestJWTExpirationMin = 10
)

func Setup() {
	DB = db.NewTestDB(TestDSN)
	RunMigrations()
}

func Teardown() {
	DB.Close()
}

func RunMigrations() {
	// Use CREATE IF NOT EXISTS to avoid races with parallel test packages.
	// Use ALTER TABLE to ensure schema evolution columns exist.

	DB.Exec(`
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
	DB.Exec(`CREATE INDEX IF NOT EXISTS idx_episodes_feed_id ON episodes(feed_id)`)
	DB.Exec(`CREATE INDEX IF NOT EXISTS idx_episodes_feed_guid ON episodes(feed_guid)`)

	DB.Exec(`
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
	DB.Exec(`CREATE INDEX IF NOT EXISTS idx_feeds_user_id ON feeds(user_id)`)

	DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255),
			google_id VARCHAR(255) UNIQUE
		)
	`)
	// Ensure columns added by later migrations exist
	DB.Exec(`ALTER TABLE users ALTER COLUMN password DROP NOT NULL`)
	DB.Exec(`ALTER TABLE users ADD COLUMN IF NOT EXISTS google_id VARCHAR(255) UNIQUE`)
	DB.Exec(`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`)
	DB.Exec(`CREATE INDEX IF NOT EXISTS idx_users_google_id ON users(google_id)`)
}

func NewApp() *echo.Echo {
	r := router.NewTestRouter()
	apiGroup := r.Group("/api")

	cfg := &config.Config{
		Auth: config.Auth{
			JwtSecret:        TestJWTSecret,
			JwtExpirationMin: TestJWTExpirationMin,
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
	// Use DELETE instead of TRUNCATE to avoid ACCESS EXCLUSIVE locks
	// that interfere with store tests running in parallel.
	// Only delete users created by integration tests (emails ending in @example.com)
	// to avoid interfering with store tests that use @test.com emails.
	DB.Exec("DELETE FROM feeds")
	DB.Exec("DELETE FROM episodes")
	DB.Exec("DELETE FROM users WHERE email LIKE '%@example.com'")
}
