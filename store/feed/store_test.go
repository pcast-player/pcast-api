package feed

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"pcast-api/db"
)

var d *sql.DB
var fs *Store

const testDSN = "host=localhost port=5432 user=pcast password=pcast dbname=pcast_test sslmode=disable"

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

	fs = New(d)
}

func tearDown() {
	// Clean up test data
	truncateTable()
	d.Close()
}

func runMigrations() {
	// Create episodes table first (from migration 00001)
	_, err := d.Exec(`
		CREATE TABLE IF NOT EXISTS episodes (
			id UUID PRIMARY KEY,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			feed_id UUID NOT NULL,
			feed_guid VARCHAR(255) NOT NULL,
			current_position INTEGER,
			played BOOLEAN NOT NULL DEFAULT FALSE
		);
		CREATE INDEX IF NOT EXISTS idx_episodes_feed_id ON episodes(feed_id);
		CREATE INDEX IF NOT EXISTS idx_episodes_feed_guid ON episodes(feed_guid);
	`)
	if err != nil {
		panic(fmt.Sprintf("failed to run episode migrations: %v", err))
	}

	// Create feeds table (from migration 00002)
	_, err = d.Exec(`
		CREATE TABLE IF NOT EXISTS feeds (
			id UUID PRIMARY KEY,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			user_id UUID NOT NULL,
			title VARCHAR(500) NOT NULL,
			url VARCHAR(1000) NOT NULL,
			synced_at TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_feeds_user_id ON feeds(user_id);
	`)
	if err != nil {
		panic(fmt.Sprintf("failed to run feed migrations: %v", err))
	}
}

func truncateTable() {
	_, err := d.Exec("TRUNCATE TABLE feeds")
	if err != nil {
		// Table might not exist yet, ignore error
		return
	}
}

func TestCreateFeed(t *testing.T) {
	userID, _ := uuid.NewV7()
	feed := &Feed{URL: "https://example.com", Title: "Example Feed", UserID: userID}
	err := fs.Create(feed)
	assert.NoError(t, err)

	truncateTable()
}

func TestFindFeedByID(t *testing.T) {
	userID, _ := uuid.NewV7()
	feed := &Feed{URL: "https://example.com", Title: "Example Feed", UserID: userID}
	err := fs.Create(feed)
	assert.NoError(t, err)

	foundFeed, err := fs.FindByID(feed.ID)

	assert.NoError(t, err)
	assert.Equal(t, feed.URL, foundFeed.URL)

	truncateTable()
}

func TestStore_FindByUserID(t *testing.T) {
	userID, err := uuid.NewV7()
	assert.NoError(t, err)
	feed := &Feed{URL: "https://example.com", Title: "Example Feed", UserID: userID}
	err = fs.Create(feed)
	assert.NoError(t, err)

	foundFeeds, err := fs.FindByUserID(userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, foundFeeds)

	truncateTable()
}

func TestStore_FindByIdAndUserID(t *testing.T) {
	userID, err := uuid.NewV7()
	assert.NoError(t, err)
	feed := &Feed{URL: "https://example.com", Title: "Example Feed", UserID: userID}
	err = fs.Create(feed)
	assert.NoError(t, err)

	foundFeed, err := fs.FindByIdAndUserID(feed.ID, userID)
	assert.NoError(t, err)
	assert.Equal(t, feed.URL, foundFeed.URL)

	truncateTable()
}

func TestDeleteFeed(t *testing.T) {
	userID, _ := uuid.NewV7()
	feed := &Feed{URL: "https://example.com", Title: "Example Feed", UserID: userID}
	err := fs.Create(feed)
	assert.NoError(t, err)

	err = fs.Delete(feed)
	assert.NoError(t, err)

	truncateTable()
}

func TestUpdateFeed(t *testing.T) {
	userID, _ := uuid.NewV7()
	feed := &Feed{URL: "https://example.com", Title: "Example Feed", UserID: userID}
	err := fs.Create(feed)
	assert.NoError(t, err)

	feed.URL = "https://example.com/updated"
	err = fs.Update(feed)
	assert.NoError(t, err)

	foundFeed, err := fs.FindByID(feed.ID)
	assert.NoError(t, err)
	assert.Equal(t, feed.URL, foundFeed.URL)

	truncateTable()
}
