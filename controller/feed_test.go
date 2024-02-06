package controller

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/steinfletcher/apitest-jsonpath"
	"gorm.io/gorm"
	"net/http"
	"os"
	"pcast-api/db"
	"pcast-api/router"
	"pcast-api/store"
	"testing"

	"github.com/steinfletcher/apitest"
)

var d *gorm.DB

func TestMain(m *testing.M) {
	d = db.NewTestDB()
	db.AutoMigrate(d)

	code := m.Run()

	db.RemoveTables(d)

	os.Exit(code)
}

func newApp() *echo.Echo {
	r := router.NewTestRouter()
	apiV1 := r.Group("/api")

	feedStore := store.New(d)
	feedController := New(feedStore)
	feedController.Register(apiV1)

	return r
}

func TestGetFeeds(t *testing.T) {
	apitest.New().
		Handler(newApp()).
		Get("/api/feeds").
		Expect(t).
		Assert(jsonpath.Len("$", 0)).
		Status(http.StatusOK).
		End()
}

func TestCreateFeed(t *testing.T) {
	apitest.New().
		Handler(newApp()).
		Post("/api/feeds").
		JSON(`{"url": "https://example.com"}`).
		Expect(t).
		Status(http.StatusCreated).
		End()

	apitest.New().
		Handler(newApp()).
		Get("/api/feeds").
		Expect(t).
		Assert(jsonpath.Len("$", 1)).
		Status(http.StatusOK).
		End()

	db.TruncateTables(d)
}

func TestCreateFeedPropertyNameError(t *testing.T) {
	apitest.New().
		Handler(newApp()).
		Post("/api/feeds").
		JSON(`{"ur": "https://example.com"}`).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	db.TruncateTables(d)
}

func TestCreateFeedUrlValidationError(t *testing.T) {
	apitest.New().
		Handler(newApp()).
		Post("/api/feeds").
		JSON(`{"url": "://example.com"}`).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	db.TruncateTables(d)
}

func TestDeleteFeed(t *testing.T) {
	apitest.New().
		Handler(newApp()).
		Post("/api/feeds").
		JSON(`{"url": "https://example.com"}`).
		Expect(t).
		Status(http.StatusCreated).
		End()

	apitest.New().
		Handler(newApp()).
		Get("/api/feeds").
		Expect(t).
		Assert(jsonpath.Len("$", 1)).
		Status(http.StatusOK).
		End()

	result, err := store.New(d).FindAll()

	if err != nil {
		t.Fatal(err)
	}

	apitest.New().
		Handler(newApp()).
		Delete(fmt.Sprintf("/api/feeds/%s", result[0].ID)).
		Expect(t).
		Status(http.StatusOK).
		End()

	apitest.New().
		Handler(newApp()).
		Get("/api/feeds").
		Expect(t).
		Assert(jsonpath.Len("$", 0)).
		Status(http.StatusOK).
		End()
}
