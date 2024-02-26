package controller

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"pcast-api/config"
	"pcast-api/controller/feed"
	"pcast-api/controller/user"
	feedService "pcast-api/service/feed"
	userService "pcast-api/service/user"
	feedStore "pcast-api/store/feed"
	userStore "pcast-api/store/user"
)

type Controller struct {
	config *config.Config
	db     *gorm.DB
}

func NewController(config *config.Config, db *gorm.DB, g *echo.Group) *Controller {
	newFeedHandler(db, g)
	newUserHandler(db, g)

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

func newUserHandler(db *gorm.DB, g *echo.Group) {
	store := userStore.New(db)
	service := userService.NewService(store)
	handler := user.NewHandler(service)

	handler.Register(g)
}
