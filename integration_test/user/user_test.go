package user

import (
	"github.com/labstack/echo/v4"
	"github.com/steinfletcher/apitest"
	"gorm.io/gorm"
	"net/http"
	"os"
	"pcast-api/controller"
	"pcast-api/db"
	"pcast-api/helper"
	"pcast-api/router"
	store "pcast-api/store/user"
	"testing"
)

var d *gorm.DB

func TestMain(m *testing.M) {
	d = db.NewTestDB("./../../fixtures/test/pcast.db")

	code := m.Run()

	helper.RemoveTable(d, &store.User{})

	os.Exit(code)
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
		Post("/api/user").
		JSON(`{"email": "foo@bar.com", "password": "test"}`).
		Expect(t).
		Status(http.StatusCreated).
		End()

	truncateTable()
}
