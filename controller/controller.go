package controller

import (
	"database/sql"

	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"pcast-api/config"
	"pcast-api/controller/feed"
	"pcast-api/controller/user"
	authMiddleware "pcast-api/middleware/auth"
	feedService "pcast-api/service/feed"
	userService "pcast-api/service/user"
	feedStore "pcast-api/store/feed"
	userStore "pcast-api/store/user"
)

type Controller struct {
	config *config.Config
	db     *sql.DB
}

// NewController initializes all handlers
func NewController(config *config.Config, db *sql.DB, g *echo.Group) *Controller {
	protected := g.Group("")
	protected.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.Auth.JwtSecret),
	}))

	newFeedHandler(db, protected, []byte(config.Auth.JwtSecret))
	newUserHandler(config, db, g, protected)

	return &Controller{
		config: config,
		db:     db,
	}
}

func newFeedHandler(db *sql.DB, g *echo.Group, jwtSecret []byte) {
	store := feedStore.New(db)
	service := feedService.NewService(store)
	middleware := authMiddleware.NewJWTMiddleware(jwtSecret)
	handler := feed.NewHandler(service, middleware)

	handler.Register(g)
}

func newUserHandler(config *config.Config, db *sql.DB, public *echo.Group, protected *echo.Group) {
	store := userStore.New(db)
	service := userService.NewService(store, config.Auth.JwtSecret, config.Auth.JwtExpirationMin)
	middleware := authMiddleware.NewJWTMiddleware([]byte(config.Auth.JwtSecret))
	handler := user.NewHandler(service, middleware)

	handler.Register(public, protected)
}
