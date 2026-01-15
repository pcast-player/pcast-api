package feed

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"net/http"
	serviceInterface "pcast-api/controller/service_interface"
	model "pcast-api/store/feed"
)

type Handler struct {
	service serviceInterface.Feed
}

func NewHandler(service serviceInterface.Feed) *Handler {
	return &Handler{service: service}
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
	userID, err := getUser(c)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	feeds, err := f.service.GetFeedsByUserID(c.Request().Context(), userID)
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
func (f *Handler) CreateFeed(c echo.Context) error {
	userID, err := getUser(c)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	r := new(CreateRequest)
	if err := c.Bind(r); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(r); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	fd := model.Feed{UserID: userID, URL: r.URL, Title: r.Title, SyncedAt: nil}

	err = f.service.CreateFeed(c.Request().Context(), &fd)
	if err != nil {
		println("store error", err.Error())
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
func (f *Handler) DeleteFeed(c echo.Context) error {
	userID, err := getUser(c)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	UUID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err = f.service.DeleteFeed(c.Request().Context(), userID, UUID); err != nil {
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
func (f *Handler) SyncFeed(c echo.Context) error {
	userID, err := getUser(c)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	feedId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	err = f.service.SyncFeed(c.Request().Context(), userID, feedId)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.NoContent(http.StatusNoContent)
}

func getUser(c echo.Context) (uuid.UUID, error) {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return uuid.Nil, errors.New("missing or invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, errors.New("invalid token claims")
	}

	sub, err := claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.Parse(sub)
}

func (f *Handler) Register(g *echo.Group) {
	g.GET("/feeds", f.GetFeeds)
	g.POST("/feeds", f.CreateFeed)
	g.PUT("/feeds/:id/sync", f.SyncFeed)
	g.DELETE("/feeds/:id", f.DeleteFeed)
}
