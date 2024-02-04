package router

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"pcast-api/validation"
)

func New() *echo.Echo {
	e := echo.New()
	e.Validator = &validation.ApiValidator{Validator: validator.New()}

	return e
}
