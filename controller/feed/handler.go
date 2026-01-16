package feed

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	serviceInterface "pcast-api/controller/service_interface"
	authMiddleware "pcast-api/middleware/auth"
	httpstatus "pcast-api/middleware/httpstatus"
	model "pcast-api/store/feed"
)

type Handler struct {
	service    serviceInterface.Feed
	middleware *authMiddleware.JWTMiddleware
}

func NewHandler(service serviceInterface.Feed, middleware *authMiddleware.JWTMiddleware) *Handler {
	return &Handler{service: service, middleware: middleware}
}

// GetFeeds godoc
// @Summary Get all feeds
// @Description Retrieve all feeds from the store
// @Tags feeds
// @Produce json
// @Param Authorization header string true "User ID"
// @Success 200 {array} Presenter
// @Router /feeds [get]
func (f *Handler) GetFeeds(c echo.Context) error {
	userID, err := f.middleware.GetUserID(c)
	if err != nil {
		return c.NoContent(httpstatus.StatusUnauthorized)
	}
	feeds, err := f.service.GetFeedsByUserID(c.Request().Context(), *userID)
	if err != nil {
		return c.NoContent(httpstatus.StatusInternalServerError)
	}

	res := lo.Map(feeds, func(item model.Feed, index int) *Presenter {
		return NewPresenter(&item)
	})

	return c.JSON(httpstatus.StatusOK, res)
}

// CreateFeed godoc
// @Summary Create a new feed
// @Description Create a new feed with the data provided in the request
// @Tags feeds
// @Accept json
// @Produce json
// @Param feed body CreateRequest true "CreateRequest data"
// @Param Authorization header string true "User ID"
// @Success 201 {object} Presenter
// @Router /feeds [post]
func (f *Handler) CreateFeed(c echo.Context) error {
	userID, err := f.middleware.GetUserID(c)
	if err != nil {
		return c.NoContent(httpstatus.StatusUnauthorized)
	}
	r := new(CreateRequest)
	if err := c.Bind(r); err != nil {
		return c.NoContent(httpstatus.StatusBadRequest)
	}
	if err := c.Validate(r); err != nil {
		return c.NoContent(httpstatus.StatusBadRequest)
	}

	fd := model.Feed{UserID: *userID, URL: r.URL, Title: r.Title, SyncedAt: nil}

	err = f.service.CreateFeed(c.Request().Context(), &fd)
	if err != nil {
		c.Logger().Error("store error", err.Error())
		return c.NoContent(httpstatus.StatusBadRequest)
	}

	res := NewPresenter(&fd)

	return c.JSON(httpstatus.StatusCreated, res)
}

// DeleteFeed godoc
// @Summary Delete a feed
// @Description Delete a feed with the given feed ID
// @Tags feeds
// @Param id path string true "ID"
// @Param Authorization header string true "User ID"
// @Success 200 "Feed deleted successfully"
// @Router /feeds/{id} [delete]
func (f *Handler) DeleteFeed(c echo.Context) error {
	userID, err := f.middleware.GetUserID(c)
	if err != nil {
		return c.NoContent(httpstatus.StatusUnauthorized)
	}
	UUID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.NoContent(httpstatus.StatusBadRequest)
	}

	if err = f.service.DeleteFeed(c.Request().Context(), *userID, UUID); err != nil {
		return c.NoContent(httpstatus.StatusInternalServerError)
	}

	return c.NoContent(httpstatus.StatusOK)
}

// SyncFeed godoc
// @Summary Sync a feed
// @Description Sync a feed with the given feed ID
// @Tags feeds
// @Param id path string true "Feed ID"
// @Param Authorization header string true "User ID"
// @Success 204 "Feed synced successfully"
// @Router /feeds/{id}/sync [put]
func (f *Handler) SyncFeed(c echo.Context) error {
	userID, err := f.middleware.GetUserID(c)
	if err != nil {
		return c.NoContent(httpstatus.StatusUnauthorized)
	}
	feedId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.NoContent(httpstatus.StatusBadRequest)
	}

	err = f.service.SyncFeed(c.Request().Context(), *userID, feedId)
	if err != nil {
		return c.NoContent(httpstatus.StatusNotFound)
	}

	return c.NoContent(httpstatus.StatusNoContent)
}

func (f *Handler) Register(g *echo.Group) {
	g.GET("/feeds", f.GetFeeds)
	g.POST("/feeds", f.CreateFeed)
	g.PUT("/feeds/:id/sync", f.SyncFeed)
	g.DELETE("/feeds/:id", f.DeleteFeed)
}
