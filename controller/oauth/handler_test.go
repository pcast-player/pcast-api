package oauth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	oauthService "pcast-api/service/oauth"
)

// mockOAuthService implements serviceInterface.OAuth for testing
type mockOAuthService struct {
	authURL       string
	authURLErr    error
	callbackToken string
	callbackErr   error
}

func (m *mockOAuthService) GetGoogleAuthURL(state string) (string, error) {
	if m.authURLErr != nil {
		return "", m.authURLErr
	}
	return m.authURL, nil
}

func (m *mockOAuthService) HandleGoogleCallback(ctx context.Context, code string) (string, error) {
	if m.callbackErr != nil {
		return "", m.callbackErr
	}
	return m.callbackToken, nil
}

func TestHandler_InitiateGoogleAuth_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/auth/google", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := &mockOAuthService{
		authURL: "https://accounts.google.com/o/oauth2/auth?client_id=test",
	}
	handler := NewHandler(mockService)

	err := handler.initiateGoogleAuth(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	assert.Equal(t, "https://accounts.google.com/o/oauth2/auth?client_id=test", rec.Header().Get("Location"))

	// Verify the state cookie was set
	cookies := rec.Result().Cookies()
	var stateCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "oauth_state" {
			stateCookie = cookie
			break
		}
	}
	assert.NotNil(t, stateCookie)
	assert.NotEmpty(t, stateCookie.Value)
	assert.True(t, stateCookie.HttpOnly)
}

func TestHandler_InitiateGoogleAuth_NotConfigured(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/auth/google", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := &mockOAuthService{
		authURLErr: oauthService.ErrGoogleNotConfigured,
	}
	handler := NewHandler(mockService)

	err := handler.initiateGoogleAuth(c)
	assert.NoError(t, err) // Handler returns nil, writes JSON error
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestHandler_HandleGoogleCallback_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?code=valid-code&state=test-state", nil)
	req.AddCookie(&http.Cookie{
		Name:  "oauth_state",
		Value: "test-state",
	})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := &mockOAuthService{
		callbackToken: "jwt-token-here",
	}
	handler := NewHandler(mockService)

	err := handler.handleGoogleCallback(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "jwt-token-here")
}

func TestHandler_HandleGoogleCallback_MissingCode(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?state=test-state", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := NewHandler(&mockOAuthService{})

	err := handler.handleGoogleCallback(c)
	assert.NoError(t, err) // Handler returns nil, writes JSON error
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "missing code or state parameter")
}

func TestHandler_HandleGoogleCallback_MissingState(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?code=valid-code", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := NewHandler(&mockOAuthService{})

	err := handler.handleGoogleCallback(c)
	assert.NoError(t, err) // Handler returns nil, writes JSON error
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "missing code or state parameter")
}

func TestHandler_HandleGoogleCallback_InvalidState(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?code=valid-code&state=wrong-state", nil)
	req.AddCookie(&http.Cookie{
		Name:  "oauth_state",
		Value: "correct-state",
	})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := NewHandler(&mockOAuthService{})

	err := handler.handleGoogleCallback(c)
	assert.NoError(t, err) // Handler returns nil, writes JSON error
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid state parameter")
}

func TestHandler_HandleGoogleCallback_MissingStateCookie(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?code=valid-code&state=test-state", nil)
	// No cookie set
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := NewHandler(&mockOAuthService{})

	err := handler.handleGoogleCallback(c)
	assert.NoError(t, err) // Handler returns nil, writes JSON error
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid state parameter")
}

func TestHandler_HandleGoogleCallback_ServiceError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?code=invalid-code&state=test-state", nil)
	req.AddCookie(&http.Cookie{
		Name:  "oauth_state",
		Value: "test-state",
	})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := &mockOAuthService{
		callbackErr: errors.New("service error"),
	}
	handler := NewHandler(mockService)

	err := handler.handleGoogleCallback(c)
	assert.NoError(t, err) // Handler returns nil, writes JSON error
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "service error")
}

func TestHandler_Register(t *testing.T) {
	e := echo.New()
	g := e.Group("/api")

	handler := NewHandler(&mockOAuthService{})
	handler.Register(g)

	// Verify routes are registered
	routes := e.Routes()
	var foundGoogleAuth, foundGoogleCallback bool
	for _, route := range routes {
		if route.Path == "/api/auth/google" && route.Method == http.MethodGet {
			foundGoogleAuth = true
		}
		if route.Path == "/api/auth/google/callback" && route.Method == http.MethodGet {
			foundGoogleCallback = true
		}
	}

	assert.True(t, foundGoogleAuth, "GET /api/auth/google route should be registered")
	assert.True(t, foundGoogleCallback, "GET /api/auth/google/callback route should be registered")
}
