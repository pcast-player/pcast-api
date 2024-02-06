package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"pcast-api/router/validator"
)

func New() *echo.Echo {
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${remote_ip} [${time_rfc3339}] \"${method} ${uri} ${protocol}\" ${status} ${bytes_out} ${user_agent}\n",
	}))
	e.Validator = validator.New()

	return e
}

func NewTestRouter() *echo.Echo {
	e := echo.New()
	e.Validator = validator.New()

	return e
}
