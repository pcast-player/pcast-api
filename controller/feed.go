package controller

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"pcast-api/model"
	"pcast-api/store"
)

type FeedController struct {
	store *store.FeedStore
}

func NewFeedController(store *store.FeedStore) *FeedController {
	return &FeedController{store: store}
}

func (c *FeedController) GetFeeds(context echo.Context) error {
	feeds, _ := c.store.FindAll()

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

	err := c.store.Create(&feed)

	if err != nil {
		return context.NoContent(http.StatusBadRequest)
	}

	return context.NoContent(http.StatusCreated)
}

func (c *FeedController) Register(group *echo.Group) {
	group.GET("/feeds", c.GetFeeds)
	group.POST("/feeds", c.CreateFeed)
}
