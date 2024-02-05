package store

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"pcast-api/model"
	"testing"
)

var db *gorm.DB

func setupDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("./../pcast-test.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	tearDown()

	os.Exit(code)
}

func setup() {
	db = setupDatabase()
	err := db.AutoMigrate(&model.Feed{})

	if err != nil {
		log.Fatal(err)
	}
}

func resetDatabase() {
	err := db.Exec("DELETE FROM feeds;").Error

	if err != nil {
		log.Fatal(err)
	}
}

func tearDown() {
	err := db.Migrator().DropTable(&model.Feed{})

	if err != nil {
		log.Fatal(err)
	}
}

func TestCreateFeed(t *testing.T) {
	feedStore := New(db)

	feed := &model.Feed{URL: "https://example.com"}
	err := feedStore.Create(feed)

	assert.NoError(t, err)

	resetDatabase()
}

func TestFindFeedByID(t *testing.T) {
	feedStore := New(db)

	feed := &model.Feed{URL: "https://example.com"}

	if err := feedStore.Create(feed); err != nil {
		log.Fatal(err)
	}

	foundFeed, err := feedStore.FindByID(feed.ID)

	assert.NoError(t, err)
	assert.Equal(t, feed.URL, foundFeed.URL)

	resetDatabase()
}

func TestDeleteFeed(t *testing.T) {
	feedStore := New(db)

	feed := &model.Feed{URL: "https://example.com"}

	if err := feedStore.Create(feed); err != nil {
		log.Fatal(err)
	}

	err := feedStore.Delete(feed)

	assert.NoError(t, err)

	resetDatabase()
}
