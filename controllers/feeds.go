package controllers

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"pcast-api/models"
)

type FeedsController struct {
	db *gorm.DB
}

func NewFeedsController(db *gorm.DB) *FeedsController {
	return &FeedsController{db: db}
}

func (c *FeedsController) GetFeeds(context echo.Context) error {
	var feeds []models.Feed

	c.db.Find(&feeds)

	return context.JSON(http.StatusOK, feeds)
}

func (c *FeedsController) CreateFeed(context echo.Context) error {
	feedRequest := new(models.CreateFeedRequest)

	if err := context.Bind(feedRequest); err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	if err := context.Validate(feedRequest); err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	feed := models.Feed{URL: feedRequest.URL}

	c.db.Create(&feed)

	return context.NoContent(http.StatusCreated)
}
