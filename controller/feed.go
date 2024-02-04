package controller

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"pcast-api/model"
)

type FeedController struct {
	db *gorm.DB
}

func NewFeedController(db *gorm.DB) *FeedController {
	return &FeedController{db: db}
}

func (c *FeedController) GetFeeds(context echo.Context) error {
	var feeds []model.Feed

	c.db.Find(&feeds)

	return context.JSON(http.StatusOK, feeds)
}

func (c *FeedController) CreateFeed(context echo.Context) error {
	feedRequest := new(model.CreateFeedRequest)

	if err := context.Bind(feedRequest); err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	if err := context.Validate(feedRequest); err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	feed := model.Feed{URL: feedRequest.URL}

	c.db.Create(&feed)

	return context.NoContent(http.StatusCreated)
}
