package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"pcast-api/config"
	"pcast-api/router/validator"
	"strings"
)

func New(c *config.Config) *echo.Echo {
	e := echo.New()

	if c.Server.Logging {
		l := getLogLevel(c)

		e.Logger.SetLevel(l)
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Format: c.Server.LogFormat}))
	}

	e.Validator = validator.New()

	return e
}

func getLogLevel(c *config.Config) log.Lvl {
	switch strings.ToLower(c.Server.LogLevel) {
	case "debug":
		return log.DEBUG
	case "info":
		return log.INFO
	case "warn":
		return log.WARN
	case "error":
		return log.ERROR
	default:
		return log.INFO
	}
}

func NewTestRouter() *echo.Echo {
	e := echo.New()
	e.Validator = validator.New()

	return e
}
