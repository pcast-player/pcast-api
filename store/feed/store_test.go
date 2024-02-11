package feed

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"log"
	"os"
	"pcast-api/db"
	"testing"
)

var d *gorm.DB
var fs *Store

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	tearDown()

	os.Exit(code)
}

func setup() {
	d = db.NewTestDB("./../../fixtures/test/pcast.db")
	fs = New(d)
	err := d.AutoMigrate(&Feed{})
	if err != nil {
		log.Fatal(err)
	}
}

func tearDown() {
	fs.RemoveTables()
}

func TestCreateFeed(t *testing.T) {
	feed := &Feed{URL: "https://example.com"}
	err := fs.Create(feed)
	assert.NoError(t, err)

	fs.TruncateTables()
}

func TestFindFeedByID(t *testing.T) {
	feed := &Feed{URL: "https://example.com"}
	if err := fs.Create(feed); err != nil {
		log.Fatal(err)
	}

	foundFeed, err := fs.FindByID(feed.ID)

	assert.NoError(t, err)
	assert.Equal(t, feed.URL, foundFeed.URL)

	fs.TruncateTables()
}

func TestDeleteFeed(t *testing.T) {
	feed := &Feed{URL: "https://example.com"}
	if err := fs.Create(feed); err != nil {
		log.Fatal(err)
	}

	err := fs.Delete(feed)
	assert.NoError(t, err)

	fs.TruncateTables()
}

func TestUpdateFeed(t *testing.T) {
	feed := &Feed{URL: "https://example.com"}
	if err := fs.Create(feed); err != nil {
		log.Fatal(err)
	}

	feed.URL = "https://example.com/updated"
	err := fs.Update(feed)
	assert.NoError(t, err)

	foundFeed, err := fs.FindByID(feed.ID)
	assert.NoError(t, err)
	assert.Equal(t, feed.URL, foundFeed.URL)

	fs.TruncateTables()
}
