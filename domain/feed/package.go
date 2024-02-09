package feed

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"pcast-api/config"
	"pcast-api/domain/feed/controller"
	"pcast-api/domain/feed/model"
	"pcast-api/domain/feed/store"
)

func New(_ *config.Config, group *echo.Group, db *gorm.DB) {
	s := store.New(db)
	c := controller.New(s)

	autoMigrateFeedModel(db)

	c.Register(group)
}

func autoMigrateFeedModel(db *gorm.DB) {
	err := db.AutoMigrate(&model.Feed{})
	if err != nil {
		panic("Failed to migrate database!")
	}
}
