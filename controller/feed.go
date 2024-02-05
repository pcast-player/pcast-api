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

func New(store *store.FeedStore) *FeedController {
	return &FeedController{store: store}
}

// GetFeeds godoc
// @Summary Get all feeds
// @Description Retrieve all feeds from the store
// @Tags feeds
// @Produce json
// @Success 200 {array} model.Feed
// @Router /feeds [get]
func (c *FeedController) GetFeeds(context echo.Context) error {
	feeds, err := c.store.FindAll()

	if err != nil {
		return context.NoContent(http.StatusInternalServerError)
	}

	return context.JSON(http.StatusOK, feeds)
}

// CreateFeed godoc
// @Summary Create a new feed
// @Description Create a new feed with the data provided in the request
// @Tags feeds
// @Accept json
// @Produce json
// @Param feed body model.CreateFeedRequest true "Feed data"
// @Success 201 {string} string "Feed created successfully"
// @Router /feeds [post]
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

	return context.JSON(http.StatusCreated, feed)
}

func (c *FeedController) Register(group *echo.Group) {
	group.GET("/feeds", c.GetFeeds)
	group.POST("/feeds", c.CreateFeed)
}
