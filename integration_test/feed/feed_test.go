package feed_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/steinfletcher/apitest"
	"github.com/steinfletcher/apitest-jsonpath"
	"pcast-api/config"
	"pcast-api/controller"
	"pcast-api/controller/feed"
	"pcast-api/controller/user"
	"pcast-api/db"
	"pcast-api/router"
)

var sqlDB *sql.DB

const testDSN = "host=localhost port=5432 user=pcast password=pcast dbname=pcast_test sslmode=disable"

func TestMain(m *testing.M) {
	sqlDB = db.NewTestDB(testDSN)

	// Clean up any leftover data from previous runs (ignore errors)
	sqlDB.Exec("TRUNCATE TABLE users CASCADE")
	sqlDB.Exec("TRUNCATE TABLE feeds")
	sqlDB.Exec("TRUNCATE TABLE episodes")

	// Run migrations to create tables
	runMigrations()

	code := m.Run()

	// Clean up (ignore errors)
	sqlDB.Exec("TRUNCATE TABLE users CASCADE")
	sqlDB.Exec("TRUNCATE TABLE feeds")
	sqlDB.Exec("TRUNCATE TABLE episodes")
	sqlDB.Close()

	os.Exit(code)
}

func runMigrations() {
	// Create all tables in order - split statements to avoid race conditions

	// Create episodes table (from migration 00001)
	sqlDB.Exec(`
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
	sqlDB.Exec(`CREATE INDEX IF NOT EXISTS idx_episodes_feed_id ON episodes(feed_id)`)
	sqlDB.Exec(`CREATE INDEX IF NOT EXISTS idx_episodes_feed_guid ON episodes(feed_guid)`)

	// Create feeds table (from migration 00002)
	sqlDB.Exec(`
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
	sqlDB.Exec(`CREATE INDEX IF NOT EXISTS idx_feeds_user_id ON feeds(user_id)`)

	// Create users table (from migration 00003)
	sqlDB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL
		)
	`)
	sqlDB.Exec(`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`)
}

func newApp() *echo.Echo {
	r := router.NewTestRouter()
	apiGroup := r.Group("/api")

	cfg := &config.Config{
		Auth: config.Auth{
			JwtSecret: "testsecret",
		},
	}

	controller.NewController(cfg, sqlDB, apiGroup)

	return r
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

func truncateTables() {
	sqlDB.Exec("TRUNCATE TABLE users CASCADE")
	sqlDB.Exec("TRUNCATE TABLE feeds")
}

func createUser(t *testing.T) (uuid.UUID, string) {
	result := apitest.New().
		Handler(newApp()).
		Post("/api/user/register").
		JSON(`{"email": "foo@bar.com", "password": "test"}`).
		Expect(t).
		Status(http.StatusCreated).
		End()

	u := unmarshal[user.Presenter](t, &result)

	loginResult := apitest.New().
		Handler(newApp()).
		Post("/api/user/login").
		JSON(`{"email": "foo@bar.com", "password": "test"}`).
		Expect(t).
		Status(http.StatusOK).
		End()

	lr := unmarshal[user.LoginResponse](t, &loginResult)

	return u.ID, lr.Token
}

func TestGetFeeds(t *testing.T) {
	_, token := createUser(t)

	apitest.New().
		Handler(newApp()).
		Get("/api/feeds").
		Header("Authorization", "Bearer "+token).
		Expect(t).
		Assert(jsonpath.Len("$", 0)).
		Status(http.StatusOK).
		End()

	truncateTables()
}

func TestCreateFeed(t *testing.T) {
	_, token := createUser(t)

	apitest.New().
		Handler(newApp()).
		Post("/api/feeds").
		Header("Authorization", "Bearer "+token).
		JSON(`{"url": "https://example.com","title":"Example"}`).
		Expect(t).
		Status(http.StatusCreated).
		End()

	apitest.New().
		Handler(newApp()).
		Get("/api/feeds").
		Header("Authorization", "Bearer "+token).
		Expect(t).
		Assert(jsonpath.Len("$", 1)).
		Status(http.StatusOK).
		End()

	truncateTables()
}

func TestCreateFeedPropertyNameError(t *testing.T) {
	_, token := createUser(t)

	apitest.New().
		Handler(newApp()).
		Post("/api/feeds").
		Header("Authorization", "Bearer "+token).
		JSON(`{"ur": "https://example.com"}`).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	truncateTables()
}

func TestCreateFeedMissingPropertyError(t *testing.T) {
	_, token := createUser(t)

	apitest.New().
		Handler(newApp()).
		Post("/api/feeds").
		Header("Authorization", "Bearer "+token).
		JSON(`{"url": "https://example.com"}`).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	truncateTables()
}

func TestCreateFeedUrlValidationError(t *testing.T) {
	_, token := createUser(t)

	apitest.New().
		Handler(newApp()).
		Post("/api/feeds").
		Header("Authorization", "Bearer "+token).
		JSON(`{"url": "://example.com"}`).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	truncateTables()
}

func TestDeleteFeed(t *testing.T) {
	_, token := createUser(t)

	result := apitest.New().
		Handler(newApp()).
		Post("/api/feeds").
		Header("Authorization", "Bearer "+token).
		JSON(`{"url": "https://example.com","title":"Example"}`).
		Expect(t).
		Status(http.StatusCreated).
		End()

	fd := unmarshal[feed.Presenter](t, &result)

	apitest.New().
		Handler(newApp()).
		Get("/api/feeds").
		Header("Authorization", "Bearer "+token).
		Expect(t).
		Assert(jsonpath.Len("$", 1)).
		Status(http.StatusOK).
		End()

	apitest.New().
		Handler(newApp()).
		Delete(fmt.Sprintf("/api/feeds/%s", fd.ID)).
		Header("Authorization", "Bearer "+token).
		Expect(t).
		Status(http.StatusOK).
		End()

	apitest.New().
		Handler(newApp()).
		Get("/api/feeds").
		Header("Authorization", "Bearer "+token).
		Expect(t).
		Assert(jsonpath.Len("$", 0)).
		Status(http.StatusOK).
		End()

	truncateTables()
}

func TestUpdateFeed(t *testing.T) {
	_, token := createUser(t)

	result := apitest.New().
		Handler(newApp()).
		Post("/api/feeds").
		Header("Authorization", "Bearer "+token).
		JSON(`{"url": "https://example.com","title":"Example"}`).
		Expect(t).
		Assert(jsonpath.Equal("$.syncedAt", nil)).
		Status(http.StatusCreated).
		End()

	fd := unmarshal[feed.Presenter](t, &result)

	apitest.New().
		Handler(newApp()).
		Put(fmt.Sprintf("/api/feeds/%s/sync", fd.ID)).
		Header("Authorization", "Bearer "+token).
		Expect(t).
		Status(http.StatusNoContent).
		End()

	apitest.New().
		Handler(newApp()).
		Get("/api/feeds").
		Header("Authorization", "Bearer "+token).
		Expect(t).
		Assert(jsonpath.NotEqual("$[0].syncedAt", nil)).
		Status(http.StatusOK).
		End()

	truncateTables()
}
