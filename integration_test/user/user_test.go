package user

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/steinfletcher/apitest"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"pcast-api/controller"
	"pcast-api/controller/user"
	"pcast-api/db"
	"pcast-api/helper"
	"pcast-api/router"
	store "pcast-api/store/user"
	"testing"
)

var d *gorm.DB

func TestMain(m *testing.M) {
	d = db.NewTestDB("./../../fixtures/test/integration_user.db")

	code := m.Run()

	helper.RemoveTable(d, &store.User{})

	os.Exit(code)
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

func newApp() *echo.Echo {
	r := router.NewTestRouter()
	apiGroup := r.Group("/api")

	controller.NewController(nil, d, apiGroup)

	return r
}

func truncateTable() {
	helper.RemoveTable(d, &store.User{})
}

func TestCreateUser(t *testing.T) {
	apitest.New().
		Handler(newApp()).
		Post("/api/user/register").
		JSON(`{"email": "foo@bar.com", "password": "test"}`).
		Expect(t).
		Status(http.StatusCreated).
		End()

	truncateTable()
}

func TestUpdatePassword(t *testing.T) {
	result := apitest.New().
		Handler(newApp()).
		Post("/api/user/register").
		JSON(`{"email": "foo@bar.com", "password": "test"}`).
		Expect(t).
		Status(http.StatusCreated).
		End()

	u := unmarshal[user.Presenter](t, &result)

	apitest.New().
		Handler(newApp()).
		Put("/api/user/password").
		Header("Authorization", u.ID.String()).
		JSON(`{"oldPassword": "test", "newPassword": "test2"}`).
		Expect(t).
		Status(http.StatusOK).
		End()

	truncateTable()
}
