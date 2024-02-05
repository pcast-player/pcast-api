package controller

import (
	"net/http"
	"net/http/httptest"
	"pcast-api/router/validator"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"pcast-api/model"
)

type MockFeedStore struct {
	mock.Mock
}

func (m *MockFeedStore) FindAll() ([]model.Feed, error) {
	args := m.Called()
	return args.Get(0).([]model.Feed), args.Error(1)
}

func (m *MockFeedStore) Create(feed *model.Feed) error {
	args := m.Called(feed)
	return args.Error(0)
}

func (m *MockFeedStore) FindByID(id uuid.UUID) (*model.Feed, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Feed), args.Error(1)
}

func (m *MockFeedStore) Delete(feed *model.Feed) error {
	args := m.Called(feed)
	return args.Error(0)
}

func TestGetFeedsReturnsFeeds(t *testing.T) {
	mockStore := new(MockFeedStore)
	mockFeed := &model.Feed{URL: "https://example.com"}
	mockStore.On("FindAll").Return([]model.Feed{*mockFeed}, nil)

	feedController := New(mockStore)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/feeds", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := feedController.GetFeeds(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, `[{"id":"00000000-0000-0000-0000-000000000000","url":"https://example.com"}]`+"\n", rec.Body.String())
	mockStore.AssertExpectations(t)
}

func TestCreateFeedReturnsCreatedFeed(t *testing.T) {
	mockStore := new(MockFeedStore)
	mockStore.On("Create", mock.AnythingOfType("*model.Feed")).Return(nil)

	feedController := New(mockStore)

	e := echo.New()
	e.Validator = validator.New()
	req := httptest.NewRequest(http.MethodPost, "/feeds", strings.NewReader(`{"url":"https://example.com"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := feedController.CreateFeed(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	mockStore.AssertExpectations(t)
}

func TestCreateFeedReturnsPropertyNameValidationError(t *testing.T) {
	mockStore := new(MockFeedStore)

	feedController := New(mockStore)

	e := echo.New()
	e.Validator = validator.New()
	req := httptest.NewRequest(http.MethodPost, "/feeds", strings.NewReader(`{"ur":"https://example.com"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := feedController.CreateFeed(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "\"code=400, message=Key: 'CreateFeedRequest.URL' Error:Field validation for 'URL' failed on the 'required' tag\"\n", rec.Body.String())
	mockStore.AssertExpectations(t)
}

func TestCreateFeedReturnsUrlValidationError(t *testing.T) {
	mockStore := new(MockFeedStore)

	feedController := New(mockStore)

	e := echo.New()
	e.Validator = validator.New()
	req := httptest.NewRequest(http.MethodPost, "/feeds", strings.NewReader(`{"url":"://example.com"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := feedController.CreateFeed(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "\"code=400, message=Key: 'CreateFeedRequest.URL' Error:Field validation for 'URL' failed on the 'url' tag\"\n", rec.Body.String())
	mockStore.AssertExpectations(t)
}

func TestDeleteFeedReturnsNoContent(t *testing.T) {
	mockStore := new(MockFeedStore)
	mockFeed := &model.Feed{URL: "https://example.com"}
	mockStore.On("FindByID", mock.AnythingOfType("uuid.UUID")).Return(mockFeed, nil)
	mockStore.On("Delete", mockFeed).Return(nil)

	feedController := New(mockStore)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/feeds/3e2076e6-1c0e-4a7e-9a6e-4b8e7f8f3e84", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("3e2076e6-1c0e-4a7e-9a6e-4b8e7f8f3e84")

	err := feedController.DeleteFeed(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockStore.AssertExpectations(t)
}
