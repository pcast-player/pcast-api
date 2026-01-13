package user

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/steinfletcher/apitest"
	"pcast-api/controller"
	"pcast-api/controller/user"
	"pcast-api/db"
	"pcast-api/router"
)

var sqlDB *sql.DB

const testDSN = "host=localhost port=5432 user=pcast password=pcast dbname=pcast_test sslmode=disable"

func TestMain(m *testing.M) {
	sqlDB = db.NewTestDBSQL(testDSN)

	code := m.Run()

	// Clean up
	sqlDB.Exec("TRUNCATE TABLE users CASCADE")
	sqlDB.Close()

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

	controller.NewController(nil, sqlDB, apiGroup)

	return r
}

func truncateTable() {
	sqlDB.Exec("TRUNCATE TABLE users CASCADE")
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
