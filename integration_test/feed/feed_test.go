package feed_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/steinfletcher/apitest"
	"github.com/steinfletcher/apitest-jsonpath"
	"pcast-api/controller/feed"
	"pcast-api/controller/user"
	testhelper "pcast-api/integration_test/testhelper"
)

func TestMain(m *testing.M) {
	testhelper.Setup()

	code := m.Run()

	testhelper.Teardown()

	os.Exit(code)
}

func newApp() *echo.Echo {
	return testhelper.NewApp()
}

func unmarshal[M any](t *testing.T, result *apitest.Result) *M {
	u, err := testhelper.UnmarshalResult[M](result.Response.Body)
	if err != nil {
		t.Fatal(err)
	}
	return u
}

func truncateTables() {
	testhelper.TruncateAll()
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
