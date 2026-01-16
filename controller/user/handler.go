package user

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	serviceInterface "pcast-api/controller/service_interface"
	authMiddleware "pcast-api/middleware/auth"
	userService "pcast-api/service/user"
)

type Handler struct {
	service    serviceInterface.User
	middleware *authMiddleware.JWTMiddleware
}

func NewHandler(service serviceInterface.User, middleware *authMiddleware.JWTMiddleware) *Handler {
	return &Handler{service: service, middleware: middleware}
}

// RegisterUser godoc
// @Summary Create a new user
// @Description Register a new user with the data provided in the request
// @Tags user
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "RegisterRequest data"
// @Success 201 {object} Presenter
// @Router /user/register [post]
func (h *Handler) registerUser(c echo.Context) error {
	userRequest := new(RegisterRequest)
	if err := c.Bind(userRequest); err != nil {
		return err
	}
	if err := c.Validate(userRequest); err != nil {
		return err
	}

	ud, err := h.service.CreateUser(c.Request().Context(), userRequest.Email, userRequest.Password)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	res := NewPresenter(ud)

	return c.JSON(http.StatusCreated, res)
}

// LoginUser godoc
// @Summary Login user
// @Description Login user with the data provided in the request
// @Tags user
// @Accept json
// @Produce json
// @Param user body LoginRequest true "LoginRequest data"
// @Success 200 {object} LoginResponse
// @Router /user/login [post]
func (h *Handler) loginUser(c echo.Context) error {
	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	token, err := h.service.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	return c.JSON(http.StatusOK, LoginResponse{Token: token})
}

// UpdatePassword godoc
// @Summary Update user password
// @Description Update user password with the data provided in the request
// @Tags user
// @Accept json
// @Produce json
// @Param Authorization header string true "User ID"
// @Param passwords body UpdatePasswordRequest true "UpdatePasswordRequest data"
// @Success 200
// @Router /user/password [put]
func (h *Handler) updatePassword(c echo.Context) error {
	userID, err := h.middleware.GetUserID(c)
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	pwRequest := new(UpdatePasswordRequest)
	if err := c.Bind(pwRequest); err != nil {
		return err
	}
	if err := c.Validate(pwRequest); err != nil {
		return err
	}

	err = h.service.UpdatePassword(c.Request().Context(), *userID, pwRequest.OldPassword, pwRequest.NewPassword)
	if err != nil {
		if errors.Is(err, userService.ErrInvalidPassword) {
			return c.NoContent(http.StatusUnauthorized)
		}
		return c.NoContent(http.StatusBadRequest)
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) Register(public *echo.Group, protected *echo.Group) {
	public.POST("/user/register", h.registerUser)
	public.POST("/user/login", h.loginUser)
	protected.PUT("/user/password", h.updatePassword)
}
