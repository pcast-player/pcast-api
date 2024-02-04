package router

import (
	"github.com/labstack/echo/v4"
)

func New() *echo.Echo {
	e := echo.New()
	e.Validator = NewValidator()

	return e
}
