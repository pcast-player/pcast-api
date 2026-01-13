package controller

import (
	"database/sql"

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
	gormDB *gorm.DB
	sqlDB  *sql.DB
}

// NewController initializes all handlers
// TODO: Once all stores migrate to sqlc, remove gormDB parameter
func NewController(config *config.Config, gormDB *gorm.DB, sqlDB *sql.DB, g *echo.Group) *Controller {
	newFeedHandler(sqlDB, g)
	newUserHandler(gormDB, g)

	return &Controller{
		config: config,
		gormDB: gormDB,
		sqlDB:  sqlDB,
	}
}

func newFeedHandler(db *sql.DB, g *echo.Group) {
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
