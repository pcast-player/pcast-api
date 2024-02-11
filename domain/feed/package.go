package feed

import (
	"github.com/labstack/echo/v4"
	"pcast-api/domain/feed/controller"
	"pcast-api/store/feed"
)

func New(g *echo.Group, fs *feed.Store) {
	c := controller.New(fs)

	c.Register(g)
}
