package user

import (
	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	service "pcast-api/service/user"
	model "pcast-api/store/user"
)

type Handler struct {
	service service.Interface
}

func NewHandler(service service.Interface) *Handler {
	return &Handler{service: service}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the data provided in the request
// @Tags user
// @Accept json
// @Produce json
// @Param user body CreateRequest true "CreateRequest data"
// @Success 201 {object} Presenter
// @Router /user [post]
func (h *Handler) createUser(c echo.Context) error {
	userRequest := new(CreateRequest)
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

	err = h.service.CreateUser(&ud)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	res := NewPresenter(&ud)

	return c.JSON(http.StatusCreated, res)
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
	r := c.Request()
	uid, err := uuid.Parse(r.Header.Get("Authorization"))
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

	user, err := h.service.GetUser(uid)
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

	err = h.service.UpdateUser(user)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	return c.NoContent(http.StatusOK)

}

func (h *Handler) Register(g *echo.Group) {
	g.POST("/user", h.createUser)
	g.PUT("/user/password", h.updatePassword)
}
