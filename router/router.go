package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"pcast-api/config"
	"pcast-api/router/validator"
	"strings"
)

var logLevels = map[string]log.Lvl{
	"debug": log.DEBUG,
	"info":  log.INFO,
	"warn":  log.WARN,
	"error": log.ERROR,
}

func New(c *config.Config) *echo.Echo {
	e := echo.New()

	if c.Server.Logging {
		ll := getLogLevel(c.Server.LogLevel)

		e.Logger.SetLevel(ll)
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Format: c.Server.LogFormat}))
	}

	e.Validator = validator.New()

	return e
}

func getLogLevel(s string) log.Lvl {
	ll, ok := logLevels[strings.ToLower(s)]
	if !ok {
		return log.INFO
	}

	return ll
}

func NewTestRouter() *echo.Echo {
	e := echo.New()
	e.Validator = validator.New()

	return e
}
