package controller

import (
	"database/sql"

	"github.com/labstack/echo/v4"
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
	db     *sql.DB
}

// NewController initializes all handlers
func NewController(config *config.Config, db *sql.DB, g *echo.Group) *Controller {
	newFeedHandler(db, g)
	newUserHandler(db, g)

	return &Controller{
		config: config,
		db:     db,
	}
}

func newFeedHandler(db *sql.DB, g *echo.Group) {
	store := feedStore.New(db)
	service := feedService.NewService(store)
	handler := feed.NewHandler(service)

	handler.Register(g)
}

func newUserHandler(db *sql.DB, g *echo.Group) {
	store := userStore.New(db)
	service := userService.NewService(store)
	handler := user.NewHandler(service)

	handler.Register(g)
}
