package controller

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"net/http"
	"pcast-api/domain/feed/presenter"
	"pcast-api/domain/feed/request"
	"pcast-api/store/feed"
	"time"
)

type Controller struct {
	store feed.Interface
}

func New(store feed.Interface) *Controller {
	return &Controller{store: store}
}

// GetFeeds godoc
// @Summary Get all feeds
// @Description Retrieve all feeds from the store
// @Tags feeds
// @Produce json
// @Success 200 {array} presenter.Feed
// @Router /feeds [get]
func (f *Controller) GetFeeds(c echo.Context) error {
	feeds, err := f.store.FindAll()
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	res := lo.Map(feeds, func(item feed.Feed, index int) *presenter.Feed {
		return presenter.New(&item)
	})

	return c.JSON(http.StatusOK, res)
}

// CreateFeed godoc
// @Summary Create a new feed
// @Description Create a new feed with the data provided in the request
// @Tags feeds
// @Accept json
// @Produce json
// @Param feed body request.Feed true "Feed data"
// @Success 201 {object} presenter.Feed
// @Router /feeds [post]
func (f *Controller) CreateFeed(c echo.Context) error {
	feedRequest := new(request.Feed)
	if err := c.Bind(feedRequest); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(feedRequest); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	fd := feed.Feed{URL: feedRequest.URL, Title: feedRequest.Title, SyncedAt: nil}

	err := f.store.Create(&fd)
	if err != nil {
		println("store error", err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	res := presenter.New(&fd)

	return c.JSON(http.StatusCreated, res)
}

// DeleteFeed godoc
// @Summary Delete a feed
// @Description Delete a feed with the given ID
// @Tags feeds
// @Param id path string true "Feed ID"
// @Success 200 "Feed deleted successfully"
// @Router /feeds/{id} [delete]
func (f *Controller) DeleteFeed(c echo.Context) error {
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

// SyncFeed godoc
// @Summary Sync a feed
// @Description Sync a feed with the given ID
// @Tags feeds
// @Param id path string true "Feed ID"
// @Success 204 "Feed synced successfully"
// @Router /feeds/{id}/sync [put]
func (f *Controller) SyncFeed(c echo.Context) error {
	UUID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	fd, err := f.store.FindByID(UUID)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	now := time.Now()
	fd.SyncedAt = &now

	err = f.store.Update(fd)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}

func (f *Controller) Register(g *echo.Group) {
	g.GET("/feeds", f.GetFeeds)
	g.POST("/feeds", f.CreateFeed)
	g.PUT("/feeds/:id/sync", f.SyncFeed)
	g.DELETE("/feeds/:id", f.DeleteFeed)
}
