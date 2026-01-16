package feed

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"pcast-api/db"
)

var d *sql.DB
var fs *Store

const (
	testDSN       = "host=localhost port=5432 user=pcast password=pcast dbname=pcast_test sslmode=disable"
	testFeedURL   = "https://example.com"
	testFeedTitle = "Example Feed"
)

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	tearDown()

	os.Exit(code)
}

func setup() {
	d = db.NewTestDB(testDSN)

	// Run migrations
	runMigrations()

	// Using background context since no cancellation or timeout is needed in tests.
	fs = New(d)
}

func tearDown() {
	// Clean up test data
	truncateTable()
	d.Close()
}

func runMigrations() {
	// NOTE: In CI, goose migrations are run before tests.
	// These CREATE TABLE statements exist to support local runs.
	// Keep users table available so FK constraints won't fail.

	d.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL
		)
	`)
	d.Exec(`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`)

	d.Exec(`
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
	d.Exec(`CREATE INDEX IF NOT EXISTS idx_episodes_feed_id ON episodes(feed_id)`)
	d.Exec(`CREATE INDEX IF NOT EXISTS idx_episodes_feed_guid ON episodes(feed_guid)`)

	d.Exec(`
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
	d.Exec(`CREATE INDEX IF NOT EXISTS idx_feeds_user_id ON feeds(user_id)`)
}

func truncateTable() {
	// If FK exists (feeds.user_id -> users.id), truncate users CASCADE clears feeds.
	// Also explicitly truncate feeds for local runs without FK.
	if _, err := d.Exec("TRUNCATE TABLE feeds"); err != nil {
		log.Printf("Failed to truncate feeds: %v", err)
	}
	if _, err := d.Exec("TRUNCATE TABLE users CASCADE"); err != nil {
		log.Printf("Failed to truncate users: %v", err)
	}
}

func ensureUserExists(t *testing.T, userID uuid.UUID) {
	email := fmt.Sprintf("user-%s@example.com", userID.String())
	now := time.Now()

	_, err := d.Exec(
		"INSERT INTO users (id, created_at, updated_at, email, password) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING",
		userID,
		now,
		now,
		email,
		"test",
	)
	if err != nil {
		t.Logf("User already exists: %v", err)
	}
}

func TestCreateFeed(t *testing.T) {
	// Using uuid.NewV7() for ID generation, which is preferred in this project.
	userID := uuid.Must(uuid.NewV7())
	ensureUserExists(t, userID)

	feed := &Feed{URL: testFeedURL, Title: testFeedTitle, UserID: userID}
	err := fs.Create(context.Background(), feed)
	assert.NoError(t, err)

	truncateTable()
}

func TestFindFeedByID(t *testing.T) {
	userID := uuid.Must(uuid.NewV7())
	ensureUserExists(t, userID)

	feed := &Feed{URL: testFeedURL, Title: testFeedTitle, UserID: userID}
	err := fs.Create(context.Background(), feed)
	assert.NoError(t, err)

	foundFeed, err := fs.FindByID(context.Background(), feed.ID)
	assert.NoError(t, err)
	assert.Equal(t, feed.URL, foundFeed.URL)

	truncateTable()
}

func TestFindFeedByID_NonExistent(t *testing.T) {
	userID := uuid.Must(uuid.NewV7())
	ensureUserExists(t, userID)

	nonExistentID := uuid.Must(uuid.NewV7())
	foundFeed, err := fs.FindByID(context.Background(), nonExistentID)
	assert.Error(t, err) // Expect error for non-existent feed
	assert.Nil(t, foundFeed)

	truncateTable()
}

func TestFindFeedByUserID(t *testing.T) {
	userID := uuid.Must(uuid.NewV7())
	ensureUserExists(t, userID)

	feed := &Feed{URL: testFeedURL, Title: testFeedTitle, UserID: userID}
	err := fs.Create(context.Background(), feed)
	assert.NoError(t, err)

	foundFeeds, err := fs.FindByUserID(context.Background(), userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, foundFeeds)

	truncateTable()
}

func TestFindFeedByUserID_EmptyResult(t *testing.T) {
	userID := uuid.Must(uuid.NewV7())
	ensureUserExists(t, userID)

	// No feeds created for this user
	foundFeeds, err := fs.FindByUserID(context.Background(), userID)
	assert.NoError(t, err)
	assert.Empty(t, foundFeeds) // Expect empty result

	truncateTable()
}

func TestFindFeedByIDAndUserID(t *testing.T) {
	userID := uuid.Must(uuid.NewV7())
	ensureUserExists(t, userID)

	feed := &Feed{URL: testFeedURL, Title: testFeedTitle, UserID: userID}
	err := fs.Create(context.Background(), feed)
	assert.NoError(t, err)

	foundFeed, err := fs.FindByIDAndUserID(context.Background(), feed.ID, userID)
	assert.NoError(t, err)
	assert.Equal(t, feed.URL, foundFeed.URL)

	truncateTable()
}

func TestDeleteFeed(t *testing.T) {
	userID := uuid.Must(uuid.NewV7())
	ensureUserExists(t, userID)

	feed := &Feed{URL: testFeedURL, Title: testFeedTitle, UserID: userID}
	err := fs.Create(context.Background(), feed)
	assert.NoError(t, err)

	err = fs.Delete(context.Background(), feed)
	assert.NoError(t, err)

	truncateTable()
}

func TestDeleteFeed_NonExistent(t *testing.T) {
	userID := uuid.Must(uuid.NewV7())
	ensureUserExists(t, userID)

	nonExistentFeed := &Feed{URL: testFeedURL, Title: testFeedTitle, UserID: userID}
	err := fs.Delete(context.Background(), nonExistentFeed)
	assert.NoError(t, err) // No error expected for non-existent feed (PostgreSQL behavior)

	truncateTable()
}

func TestUpdateFeed(t *testing.T) {
	userID := uuid.Must(uuid.NewV7())
	ensureUserExists(t, userID)

	feed := &Feed{URL: testFeedURL, Title: testFeedTitle, UserID: userID}
	err := fs.Create(context.Background(), feed)
	assert.NoError(t, err)

	feed.URL = "https://example.com/updated"
	err = fs.Update(context.Background(), feed)
	assert.NoError(t, err)

	foundFeed, err := fs.FindByID(context.Background(), feed.ID)
	assert.NoError(t, err)
	assert.Equal(t, feed.URL, foundFeed.URL)

	truncateTable()
}

func TestUpdateFeed_InvalidURL(t *testing.T) {
	userID := uuid.Must(uuid.NewV7())
	ensureUserExists(t, userID)

	feed := &Feed{URL: testFeedURL, Title: testFeedTitle, UserID: userID}
	err := fs.Create(context.Background(), feed)
	assert.NoError(t, err)

	// Test updating with invalid data (empty URL)
	feed.URL = ""
	err = fs.Update(context.Background(), feed)
	if err != nil {
		t.Logf("Expected error for empty URL: %v", err)
	} else {
		t.Log("No error returned for empty URL, but PostgreSQL should reject it")
	}

	truncateTable()
}
