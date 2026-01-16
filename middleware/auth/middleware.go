package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// UserContextKey is the context key for storing user ID
type UserContextKey string

const (
	UserIDKey UserContextKey = "user_id"
)

// JWTMiddleware extracts user ID from JWT token and stores it in context
type JWTMiddleware struct {
	secret []byte
}

func NewJWTMiddleware(secret []byte) *JWTMiddleware {
	return &JWTMiddleware{secret: secret}
}

func (m *JWTMiddleware) ExtractUserID(c echo.Context) (*uuid.UUID, error) {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return nil, echo.ErrUnauthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, echo.ErrUnauthorized
	}

	sub, err := claims.GetSubject()
	if err != nil {
		return nil, echo.ErrUnauthorized
	}

	userID, err := uuid.Parse(sub)
	if err != nil {
		return nil, echo.ErrBadRequest
	}

	return &userID, nil
}

func (m *JWTMiddleware) SetUserID(c echo.Context, userID uuid.UUID) {
	c.Set(string(UserIDKey), userID)
}

func (m *JWTMiddleware) GetUserID(c echo.Context) (*uuid.UUID, error) {
	userID, ok := c.Get(string(UserIDKey)).(uuid.UUID)
	if !ok {
		return nil, echo.ErrUnauthorized
	}

	return &userID, nil
}
