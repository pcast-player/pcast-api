package store

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"log"
	"os"
	"pcast-api/db"
	"pcast-api/model"
	"testing"
)

var d *gorm.DB

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	tearDown()

	os.Exit(code)
}

func setup() {
	d = db.NewTestDB()
	err := d.AutoMigrate(&model.Feed{})

	if err != nil {
		log.Fatal(err)
	}
}

func tearDown() {
	db.RemoveTables(d)
}

func TestCreateFeed(t *testing.T) {
	feedStore := New(d)

	feed := &model.Feed{URL: "https://example.com"}
	err := feedStore.Create(feed)

	assert.NoError(t, err)

	db.TruncateTables(d)
}

func TestFindFeedByID(t *testing.T) {
	feedStore := New(d)

	feed := &model.Feed{URL: "https://example.com"}

	if err := feedStore.Create(feed); err != nil {
		log.Fatal(err)
	}

	foundFeed, err := feedStore.FindByID(feed.ID)

	assert.NoError(t, err)
	assert.Equal(t, feed.URL, foundFeed.URL)

	db.TruncateTables(d)
}

func TestDeleteFeed(t *testing.T) {
	feedStore := New(d)

	feed := &model.Feed{URL: "https://example.com"}

	if err := feedStore.Create(feed); err != nil {
		log.Fatal(err)
	}

	err := feedStore.Delete(feed)

	assert.NoError(t, err)

	db.TruncateTables(d)
}
