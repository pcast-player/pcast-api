package feed_test

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/steinfletcher/apitest-jsonpath"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"pcast-api/controller"
	"pcast-api/controller/feed"
	"pcast-api/controller/user"
	"pcast-api/db"
	"pcast-api/helper"
	"pcast-api/router"
	feedStore "pcast-api/store/feed"
	userStore "pcast-api/store/user"
	"testing"

	"github.com/steinfletcher/apitest"
)

var d *gorm.DB

func TestMain(m *testing.M) {
	d = db.NewTestDB("./../../fixtures/test/integration_feed.db")

	code := m.Run()

	helper.RemoveTable(d, &feedStore.Feed{})
	helper.RemoveTable(d, &userStore.User{})

	os.Exit(code)
}

func newApp() *echo.Echo {
	r := router.NewTestRouter()
	apiGroup := r.Group("/api")

	controller.NewController(nil, d, apiGroup)

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
	helper.TruncateTables(d, "feeds")
	helper.TruncateTables(d, "users")
}

func createUser(t *testing.T) uuid.UUID {
	result := apitest.New().
		Handler(newApp()).
		Post("/api/user").
		JSON(`{"email": "foo@bar.com", "password": "test"}`).
		Expect(t).
		Status(http.StatusCreated).
		End()

	u := unmarshal[user.Presenter](t, &result)

	return u.ID
}

func TestGetFeeds(t *testing.T) {
	userID := createUser(t)

	apitest.New().
		Handler(newApp()).
		Get("/api/feeds").
		Header("Authorization", userID.String()).
		Expect(t).
		Assert(jsonpath.Len("$", 0)).
		Status(http.StatusOK).
		End()

	truncateTables()
}

func TestCreateFeed(t *testing.T) {
	userID := createUser(t)

	apitest.New().
		Handler(newApp()).
		Post("/api/feeds").
		Header("Authorization", userID.String()).
		JSON(`{"url": "https://example.com","title":"Example"}`).
		Expect(t).
		Status(http.StatusCreated).
		End()

	apitest.New().
		Handler(newApp()).
		Get("/api/feeds").
		Header("Authorization", userID.String()).
		Expect(t).
		Assert(jsonpath.Len("$", 1)).
		Status(http.StatusOK).
		End()

	truncateTables()
}

func TestCreateFeedPropertyNameError(t *testing.T) {
	userID := createUser(t)

	apitest.New().
		Handler(newApp()).
		Post("/api/feeds").
		Header("Authorization", userID.String()).
		JSON(`{"ur": "https://example.com"}`).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	truncateTables()
}

func TestCreateFeedMissingPropertyError(t *testing.T) {
	userID := createUser(t)

	apitest.New().
		Handler(newApp()).
		Post("/api/feeds").
		Header("Authorization", userID.String()).
		JSON(`{"url": "https://example.com"}`).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	truncateTables()
}

func TestCreateFeedUrlValidationError(t *testing.T) {
	userID := createUser(t)

	apitest.New().
		Handler(newApp()).
		Post("/api/feeds").
		Header("Authorization", userID.String()).
		JSON(`{"url": "://example.com"}`).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	truncateTables()
}

func TestDeleteFeed(t *testing.T) {
	userID := createUser(t)

	result := apitest.New().
		Handler(newApp()).
		Post("/api/feeds").
		Header("Authorization", userID.String()).
		JSON(`{"url": "https://example.com","title":"Example"}`).
		Expect(t).
		Status(http.StatusCreated).
		End()

	fd := unmarshal[feed.Presenter](t, &result)

	apitest.New().
		Handler(newApp()).
		Get("/api/feeds").
		Header("Authorization", userID.String()).
		Expect(t).
		Assert(jsonpath.Len("$", 1)).
		Status(http.StatusOK).
		End()

	apitest.New().
		Handler(newApp()).
		Delete(fmt.Sprintf("/api/feeds/%s", fd.ID)).
		Header("Authorization", userID.String()).
		Expect(t).
		Status(http.StatusOK).
		End()

	apitest.New().
		Handler(newApp()).
		Get("/api/feeds").
		Header("Authorization", userID.String()).
		Expect(t).
		Assert(jsonpath.Len("$", 0)).
		Status(http.StatusOK).
		End()

	truncateTables()
}

func TestUpdateFeed(t *testing.T) {
	userID := createUser(t)

	result := apitest.New().
		Handler(newApp()).
		Post("/api/feeds").
		Header("Authorization", userID.String()).
		JSON(`{"url": "https://example.com","title":"Example"}`).
		Expect(t).
		Assert(jsonpath.Equal("$.syncedAt", nil)).
		Status(http.StatusCreated).
		End()

	fd := unmarshal[feed.Presenter](t, &result)

	apitest.New().
		Handler(newApp()).
		Put(fmt.Sprintf("/api/feeds/%s/sync", fd.ID)).
		Header("Authorization", userID.String()).
		Expect(t).
		Status(http.StatusNoContent).
		End()

	apitest.New().
		Handler(newApp()).
		Get("/api/feeds").
		Header("Authorization", userID.String()).
		Expect(t).
		Assert(jsonpath.NotEqual("$[0].syncedAt", nil)).
		Status(http.StatusOK).
		End()

	truncateTables()
}
