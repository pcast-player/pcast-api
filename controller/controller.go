package controller

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"pcast-api/config"
	"pcast-api/controller/feed"
	feedService "pcast-api/service/feed"
	feedStore "pcast-api/store/feed"
)

type Controller struct {
	config *config.Config
	db     *gorm.DB
}

func NewController(config *config.Config, db *gorm.DB, g *echo.Group) *Controller {
	newFeedHandler(db, g)

	return &Controller{
		config: config,
		db:     db,
	}
}

func newFeedHandler(db *gorm.DB, g *echo.Group) {
	store := feedStore.New(db)
	service := feedService.NewService(store)
	handler := feed.NewHandler(service)

	handler.Register(g)
}
