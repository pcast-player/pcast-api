package oauth

import (
	"net/http"

	"github.com/labstack/echo/v4"

	serviceInterface "pcast-api/controller/service_interface"
	oauthService "pcast-api/service/oauth"
)

type Handler struct {
	service serviceInterface.OAuth
}

func NewHandler(service serviceInterface.OAuth) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(public *echo.Group) {
	public.GET("/auth/google", h.initiateGoogleAuth)
	public.GET("/auth/google/callback", h.handleGoogleCallback)
}

// initiateGoogleAuth godoc
// @Summary Initiate Google OAuth
// @Description Redirects to Google OAuth consent screen
// @Tags auth
// @Produce json
// @Success 302 {string} string "Redirect to Google"
// @Failure 500 {object} map[string]string
// @Router /auth/google [get]
func (h *Handler) initiateGoogleAuth(c echo.Context) error {
	state, err := oauthService.GenerateState()
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	// Store state in a cookie for validation on callback
	c.SetCookie(&http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   c.Request().TLS != nil, // Only secure in HTTPS
		MaxAge:   600,                    // 10 minutes
		SameSite: http.SameSiteLaxMode,
	})

	url, err := h.service.GetGoogleAuthURL(state)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.Redirect(http.StatusTemporaryRedirect, url)
}

// handleGoogleCallback godoc
// @Summary Google OAuth callback
// @Description Handles Google OAuth callback and returns JWT token
// @Tags auth
// @Produce json
// @Param code query string true "Authorization code from Google"
// @Param state query string true "State parameter for CSRF validation"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/google/callback [get]
func (h *Handler) handleGoogleCallback(c echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")

	if code == "" || state == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "missing code or state parameter",
		})
	}

	// Validate state from cookie
	cookie, err := c.Cookie("oauth_state")
	if err != nil || cookie.Value != state {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid state parameter",
		})
	}

	// Clear the state cookie
	c.SetCookie(&http.Cookie{
		Name:     "oauth_state",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	token, err := h.service.HandleGoogleCallback(c.Request().Context(), code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, LoginResponse{Token: token})
}
