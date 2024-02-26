package user

import (
	"github.com/alexedwards/argon2id"
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

	return c.NoContent(http.StatusCreated)
}

func (h *Handler) Register(g *echo.Group) {
	g.POST("/user", h.createUser)
}
