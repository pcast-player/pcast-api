package controller

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"pcast-api/feed"
	"pcast-api/model"
)

type FeedController struct {
	store feed.Store
}

func New(store feed.Store) *FeedController {
	return &FeedController{store: store}
}

// GetFeeds godoc
// @Summary Get all feeds
// @Description Retrieve all feeds from the store
// @Tags feeds
// @Produce json
// @Success 200 {array} model.Feed
// @Router /feeds [get]
func (f *FeedController) GetFeeds(c echo.Context) error {
	feeds, err := f.store.FindAll()
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, feeds)
}

// CreateFeed godoc
// @Summary Create a new feed
// @Description Create a new feed with the data provided in the request
// @Tags feeds
// @Accept json
// @Produce json
// @Param feed body model.CreateFeedRequest true "Feed data"
// @Success 201 {object} model.Feed
// @Router /feeds [post]
func (f *FeedController) CreateFeed(c echo.Context) error {
	feedRequest := new(model.CreateFeedRequest)
	if err := c.Bind(feedRequest); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(feedRequest); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	fd := model.Feed{URL: feedRequest.URL}

	err := f.store.Create(&fd)
	if err != nil {
		println("store error", err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	return c.JSON(http.StatusCreated, fd)
}

// DeleteFeed godoc
// @Summary Delete a feed
// @Description Delete a feed with the given ID
// @Tags feeds
// @Param id path string true "Feed ID"
// @Success 200 "Feed deleted successfully"
// @Router /feeds/{id} [delete]
func (f *FeedController) DeleteFeed(c echo.Context) error {
	UUID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	fd, err := f.store.FindByID(UUID)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	if err = f.store.Delete(fd); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func (f *FeedController) Register(g *echo.Group) {
	g.GET("/feeds", f.GetFeeds)
	g.POST("/feeds", f.CreateFeed)
	g.DELETE("/feeds/:id", f.DeleteFeed)
}
