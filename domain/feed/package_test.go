package feed

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/steinfletcher/apitest-jsonpath"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"pcast-api/config"
	"pcast-api/db"
	"pcast-api/domain/feed/model"
	"pcast-api/router"
	"testing"

	"github.com/steinfletcher/apitest"
)

var d *gorm.DB

func TestMain(m *testing.M) {
	d = db.NewTestDB("./../../fixtures/test/pcast.db")

	code := m.Run()

	db.RemoveTables(d)

	os.Exit(code)
}

func newApp() *echo.Echo {
	c := config.New("./../../fixtures/test/config.toml")
	r := router.NewTestRouter()
	apiV1 := r.Group("/api")

	New(c, apiV1, d)

	return r
}

func unmarshal[M any](t *testing.T, result *apitest.Result) *M {
	bytes, err := io.ReadAll(result.Response.Body)
	if err != nil {
		t.Fatal(err)
	}

	body := string(bytes)
	println(fmt.Sprintf("Response: %s", body))
	m := new(M)
	err = json.Unmarshal([]byte(body), m)
	if err != nil {
		t.Fatal(err)
	}

	return m
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
		JSON(`{"url": "https://example.com","title":"Example"}`).
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

func TestCreateFeedMissingPropertyError(t *testing.T) {
	apitest.New().
		Handler(newApp()).
		Post("/api/feeds").
		JSON(`{"url": "https://example.com"}`).
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
		JSON(`{"url": "https://example.com","title":"Example"}`).
		Expect(t).
		Status(http.StatusCreated).
		End()

	fd := unmarshal[model.Feed](t, &result)

	apitest.New().
		Handler(newApp()).
		Get("/api/feeds").
		Expect(t).
		Assert(jsonpath.Len("$", 1)).
		Status(http.StatusOK).
		End()

	apitest.New().
		Handler(newApp()).
		Delete(fmt.Sprintf("/api/feeds/%s", fd.ID)).
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

	db.TruncateTables(d)
}

func TestUpdateFeed(t *testing.T) {
	result := apitest.New().
		Handler(newApp()).
		Post("/api/feeds").
		JSON(`{"url": "https://example.com","title":"Example"}`).
		Expect(t).
		Assert(jsonpath.Equal("$.syncedAt", nil)).
		Status(http.StatusCreated).
		End()

	fd := unmarshal[model.Feed](t, &result)

	apitest.New().
		Handler(newApp()).
		Put(fmt.Sprintf("/api/feeds/%s/sync", fd.ID)).
		Expect(t).
		Status(http.StatusNoContent).
		End()

	apitest.New().
		Handler(newApp()).
		Get("/api/feeds").
		Expect(t).
		Assert(jsonpath.NotEqual("$[0].syncedAt", nil)).
		Status(http.StatusOK).
		End()

	db.TruncateTables(d)
}
