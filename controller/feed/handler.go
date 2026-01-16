package feed

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"

	serviceInterface "pcast-api/controller/service_interface"
	authMiddleware "pcast-api/middleware/auth"
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
func (h *Handler) GetFeeds(c echo.Context) error {
	userID, err := h.middleware.GetUserID(c)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	feeds, err := h.service.GetFeedsByUserID(c.Request().Context(), *userID)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	res := lo.Map(feeds, func(item model.Feed, index int) *Presenter {
		return NewPresenter(&item)
	})

	return c.JSON(http.StatusOK, res)
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
func (h *Handler) CreateFeed(c echo.Context) error {
	userID, err := h.middleware.GetUserID(c)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	r := new(CreateRequest)
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}

	fd := model.Feed{UserID: *userID, URL: r.URL, Title: r.Title, SyncedAt: nil}

	err = h.service.CreateFeed(c.Request().Context(), &fd)
	if err != nil {
		c.Logger().Error("store error", err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	res := NewPresenter(&fd)

	return c.JSON(http.StatusCreated, res)
}

// DeleteFeed godoc
// @Summary Delete a feed
// @Description Delete a feed with the given feed ID
// @Tags feeds
// @Param id path string true "ID"
// @Param Authorization header string true "User ID"
// @Success 200 "Feed deleted successfully"
// @Router /feeds/{id} [delete]
func (h *Handler) DeleteFeed(c echo.Context) error {
	userID, err := h.middleware.GetUserID(c)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	feedID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err = h.service.DeleteFeed(c.Request().Context(), *userID, feedID); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

// SyncFeed godoc
// @Summary Sync a feed
// @Description Sync a feed with the given feed ID
// @Tags feeds
// @Param id path string true "Feed ID"
// @Param Authorization header string true "User ID"
// @Success 204 "Feed synced successfully"
// @Router /feeds/{id}/sync [put]
func (h *Handler) SyncFeed(c echo.Context) error {
	userID, err := h.middleware.GetUserID(c)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	feedID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	err = h.service.SyncFeed(c.Request().Context(), *userID, feedID)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) Register(g *echo.Group) {
	g.GET("/feeds", h.GetFeeds)
	g.POST("/feeds", h.CreateFeed)
	g.PUT("/feeds/:id/sync", h.SyncFeed)
	g.DELETE("/feeds/:id", h.DeleteFeed)
}
