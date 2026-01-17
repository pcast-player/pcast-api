package controller

import (
	"database/sql"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"

	"pcast-api/config"
	"pcast-api/controller/feed"
	"pcast-api/controller/oauth"
	"pcast-api/controller/user"
	authMiddleware "pcast-api/middleware/auth"
	feedService "pcast-api/service/feed"
	oauthService "pcast-api/service/oauth"
	userService "pcast-api/service/user"
	feedStore "pcast-api/store/feed"
	userStore "pcast-api/store/user"
)

// NewController initializes all handlers
func NewController(config *config.Config, db *sql.DB, g *echo.Group) {
	middleware := authMiddleware.NewJWTMiddleware([]byte(config.Auth.JwtSecret))

	protected := g.Group("")
	protected.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.Auth.JwtSecret),
	}))
	protected.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID, err := middleware.ExtractUserID(c)
			if err != nil {
				return err
			}
			middleware.SetUserID(c, *userID)
			return next(c)
		}
	})

	newFeedHandler(db, protected, middleware)
	newUserHandler(config, db, g, protected, middleware)
	newOAuthHandler(config, db, g)
}

func newFeedHandler(db *sql.DB, g *echo.Group, middleware *authMiddleware.JWTMiddleware) {
	store := feedStore.New(db)
	service := feedService.NewService(store)
	handler := feed.NewHandler(service, middleware)

	handler.Register(g)
}

func newUserHandler(config *config.Config, db *sql.DB, public *echo.Group, protected *echo.Group, middleware *authMiddleware.JWTMiddleware) {
	store := userStore.New(db)
	service := userService.NewService(store, config.Auth.JwtSecret, config.Auth.JwtExpirationMin)
	handler := user.NewHandler(service, middleware)

	handler.Register(public, protected)
}

func newOAuthHandler(config *config.Config, db *sql.DB, public *echo.Group) {
	store := userStore.New(db)
	service := oauthService.NewService(config, store)
	handler := oauth.NewHandler(service)

	handler.Register(public)
}
