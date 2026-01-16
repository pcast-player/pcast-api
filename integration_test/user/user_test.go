package user

import (
	"net/http"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/steinfletcher/apitest"
	"pcast-api/controller/user"
	testhelper "pcast-api/integration_test/testhelper"
)

func TestMain(m *testing.M) {
	testhelper.Setup()

	code := m.Run()

	testhelper.Teardown()

	os.Exit(code)
}

func unmarshal[M any](t *testing.T, result *apitest.Result) *M {
	u, err := testhelper.UnmarshalResult[M](result.Response.Body)
	if err != nil {
		t.Fatal(err)
	}
	return u
}

func newApp() *echo.Echo {
	return testhelper.NewApp()
}

func truncateTable() {
	testhelper.TruncateAll()
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
	apitest.New().
		Handler(newApp()).
		Post("/api/user/register").
		JSON(`{"email": "foo@bar.com", "password": "test"}`).
		Expect(t).
		Status(http.StatusCreated).
		End()

	// Login to get token
	loginResult := apitest.New().
		Handler(newApp()).
		Post("/api/user/login").
		JSON(`{"email": "foo@bar.com", "password": "test"}`).
		Expect(t).
		Status(http.StatusOK).
		End()

	lr := unmarshal[user.LoginResponse](t, &loginResult)

	apitest.New().
		Handler(newApp()).
		Put("/api/user/password").
		Header("Authorization", "Bearer "+lr.Token).
		JSON(`{"oldPassword": "test", "newPassword": "test2"}`).
		Expect(t).
		Status(http.StatusOK).
		End()

	truncateTable()
}
