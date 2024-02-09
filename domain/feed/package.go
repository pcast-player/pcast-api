package feed

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"pcast-api/config"
	"pcast-api/domain/feed/controller"
	"pcast-api/domain/feed/store"
)

func New(config *config.Config, group *echo.Group, db *gorm.DB) {
	s := store.New(db)
	c := controller.New(s)

	c.Register(group)
}
