package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/steinfletcher/apitest"
	"pcast-api/controller"
	"pcast-api/controller/user"
	"pcast-api/db"
	"pcast-api/router"
)

var sqlDB *sql.DB

const testDSN = "host=localhost port=5432 user=pcast password=pcast dbname=pcast_test sslmode=disable"

func TestMain(m *testing.M) {
	sqlDB = db.NewTestDBSQL(testDSN)

	// Clean up any leftover data from previous runs
	sqlDB.Exec("TRUNCATE TABLE IF EXISTS users CASCADE")
	sqlDB.Exec("TRUNCATE TABLE IF EXISTS feeds CASCADE")
	sqlDB.Exec("TRUNCATE TABLE IF EXISTS episodes CASCADE")

	// Run migrations to create tables
	runMigrations()

	code := m.Run()

	// Clean up
	sqlDB.Exec("TRUNCATE TABLE IF EXISTS users CASCADE")
	sqlDB.Exec("TRUNCATE TABLE IF EXISTS feeds CASCADE")
	sqlDB.Exec("TRUNCATE TABLE IF EXISTS episodes CASCADE")
	sqlDB.Close()

	os.Exit(code)
}

func runMigrations() {
	// Create all tables in order - split statements to avoid race conditions

	// Create episodes table (from migration 00001)
	_, err := sqlDB.Exec(`
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
		fmt.Printf("Warning: episodes table creation: %v\n", err)
	}
	sqlDB.Exec(`CREATE INDEX IF NOT EXISTS idx_episodes_feed_id ON episodes(feed_id)`)
	sqlDB.Exec(`CREATE INDEX IF NOT EXISTS idx_episodes_feed_guid ON episodes(feed_guid)`)

	// Create feeds table (from migration 00002)
	_, err = sqlDB.Exec(`
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
		fmt.Printf("Warning: feeds table creation: %v\n", err)
	}
	sqlDB.Exec(`CREATE INDEX IF NOT EXISTS idx_feeds_user_id ON feeds(user_id)`)

	// Create users table (from migration 00003)
	_, err = sqlDB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL
		)
	`)
	if err != nil {
		panic(fmt.Sprintf("CRITICAL: failed to create users table: %v", err))
	}
	sqlDB.Exec(`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`)
}

func unmarshal[M any](t *testing.T, result *apitest.Result) *M {
	bytes, err := io.ReadAll(result.Response.Body)
	if err != nil {
		t.Fatal(err)
	}

	body := string(bytes)
	m := new(M)
	err = json.Unmarshal([]byte(body), m)
	if err != nil {
		t.Fatal(err)
	}

	return m
}

func newApp() *echo.Echo {
	r := router.NewTestRouter()
	apiGroup := r.Group("/api")

	controller.NewController(nil, sqlDB, apiGroup)

	return r
}

func truncateTable() {
	sqlDB.Exec("TRUNCATE TABLE users CASCADE")
}

func TestCreateUser(t *testing.T) {
	apitest.New().
		Handler(newApp()).
		Post("/api/user/register").
		JSON(`{"email": "foo@bar.com", "password": "test"}`).
		Expect(t).
		Status(http.StatusCreated).
		End()

	truncateTable()
}

func TestUpdatePassword(t *testing.T) {
	result := apitest.New().
		Handler(newApp()).
		Post("/api/user/register").
		JSON(`{"email": "foo@bar.com", "password": "test"}`).
		Expect(t).
		Status(http.StatusCreated).
		End()

	u := unmarshal[user.Presenter](t, &result)

	apitest.New().
		Handler(newApp()).
		Put("/api/user/password").
		Header("Authorization", u.ID.String()).
		JSON(`{"oldPassword": "test", "newPassword": "test2"}`).
		Expect(t).
		Status(http.StatusOK).
		End()

	truncateTable()
}
