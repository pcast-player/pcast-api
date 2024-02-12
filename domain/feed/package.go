package feed

import (
	"github.com/labstack/echo/v4"
	"pcast-api/store/feed"
)

func New(g *echo.Group, fs *feed.Store) {
	c := NewController(fs)

	c.Register(g)
}
