package controller

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/steinfletcher/apitest-jsonpath"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"pcast-api/db"
	"pcast-api/model"
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

func getBody(t *testing.T, result *apitest.Result) string {
	bytes, err := io.ReadAll(result.Response.Body)

	if err != nil {
		t.Fatal(err)
	}

	return string(bytes)
}

func unmarshal(t *testing.T, result *apitest.Result, v interface{}) {
	body := getBody(t, result)

	err := json.Unmarshal([]byte(body), v)

	if err != nil {
		t.Fatal(err)
	}
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
	result := apitest.New().
		Handler(newApp()).
		Post("/api/feeds").
		JSON(`{"url": "https://example.com"}`).
		Expect(t).
		Status(http.StatusCreated).
		End()

	var feed model.Feed
	unmarshal(t, &result, &feed)

	apitest.New().
		Handler(newApp()).
		Get("/api/feeds").
		Expect(t).
		Assert(jsonpath.Len("$", 1)).
		Status(http.StatusOK).
		End()

	apitest.New().
		Handler(newApp()).
		Delete(fmt.Sprintf("/api/feeds/%s", feed.ID)).
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
