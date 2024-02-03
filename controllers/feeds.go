package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"pcast-api/models"
)

func GetFeeds(c echo.Context) error {
	return c.JSON(http.StatusOK, models.GetFeeds())
}

func CreateFeed(c echo.Context) error {
	feedRequest := new(models.CreateFeedRequest)

	if err := c.Bind(feedRequest); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(feedRequest); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	feed := models.Feed{URL: feedRequest.URL}

	models.CreateFeed(&feed)

	return c.NoContent(http.StatusCreated)
}
