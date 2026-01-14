package user

import (
	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	serviceInterface "pcast-api/controller/service_interface"
	model "pcast-api/store/user"
)

type Handler struct {
	service serviceInterface.User
}

func NewHandler(service serviceInterface.User) *Handler {
	return &Handler{service: service}
}

// RegisterUser godoc
// @Summary Create a new user
// @Description Register a new user with the data provided in the request
// @Tags user
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "RegisterRequest data"
// @Success 201 {object} Presenter
// @Router /user [post]
func (h *Handler) registerUser(c echo.Context) error {
	userRequest := new(RegisterRequest)
	if err := c.Bind(userRequest); err != nil {
		return err
	}
	if err := c.Validate(userRequest); err != nil {
		return err
	}

	hash, err := argon2id.CreateHash(userRequest.Password, argon2id.DefaultParams)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	ud := model.User{Email: userRequest.Email, Password: hash}

	err = h.service.CreateUser(c.Request().Context(), &ud)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	res := NewPresenter(&ud)

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
		return c.JSON(http.StatusUnauthorized, err.Error())
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
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.NoContent(http.StatusUnauthorized)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.NoContent(http.StatusUnauthorized)
	}

	sub, err := claims.GetSubject()
	if err != nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	uid, err := uuid.Parse(sub)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	pwRequest := new(UpdatePasswordRequest)
	if err := c.Bind(pwRequest); err != nil {
		return err
	}
	if err := c.Validate(pwRequest); err != nil {
		return err
	}

	user, err := h.service.GetUser(c.Request().Context(), uid)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)

	}

	match, err := argon2id.ComparePasswordAndHash(pwRequest.OldPassword, user.Password)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	if !match {
		return c.NoContent(http.StatusUnauthorized)
	}

	hash, err := argon2id.CreateHash(pwRequest.NewPassword, argon2id.DefaultParams)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	user.Password = hash

	err = h.service.UpdateUser(c.Request().Context(), user)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	return c.NoContent(http.StatusOK)

}

func (h *Handler) Register(public *echo.Group, protected *echo.Group) {
	public.POST("/user/register", h.registerUser)
	public.POST("/user/login", h.loginUser)
	protected.PUT("/user/password", h.updatePassword)
}
