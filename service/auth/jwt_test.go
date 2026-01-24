package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateJWTToken(t *testing.T) {
	userID := uuid.Must(uuid.NewV7())
	secret := "testsecret"
	expirationMin := 60

	tokenString, err := CreateJWTToken(userID, secret, expirationMin)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Parse the token and verify claims
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	assert.NoError(t, err)
	assert.True(t, token.Valid)

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	assert.True(t, ok)
	assert.Equal(t, userID.String(), claims.Subject)
	assert.NotNil(t, claims.ExpiresAt)

	// Verify expiration is approximately correct (within a second)
	expectedExpiration := time.Now().Add(time.Duration(expirationMin) * time.Minute)
	assert.WithinDuration(t, expectedExpiration, claims.ExpiresAt.Time, 2*time.Second)
}

func TestCreateJWTToken_DifferentExpiration(t *testing.T) {
	userID := uuid.Must(uuid.NewV7())
	secret := "testsecret"

	// Test with different expiration times
	testCases := []int{1, 10, 120, 1440}

	for _, expirationMin := range testCases {
		tokenString, err := CreateJWTToken(userID, secret, expirationMin)
		assert.NoError(t, err)
		assert.NotEmpty(t, tokenString)

		// Parse and verify expiration
		token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		assert.NoError(t, err)

		claims := token.Claims.(*jwt.RegisteredClaims)
		expectedExpiration := time.Now().Add(time.Duration(expirationMin) * time.Minute)
		assert.WithinDuration(t, expectedExpiration, claims.ExpiresAt.Time, 2*time.Second)
	}
}

func TestCreateJWTToken_SigningMethod(t *testing.T) {
	userID := uuid.Must(uuid.NewV7())
	secret := "testsecret"

	tokenString, err := CreateJWTToken(userID, secret, 60)
	assert.NoError(t, err)

	// Parse without validation to check the signing method
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, &jwt.RegisteredClaims{})
	assert.NoError(t, err)
	assert.Equal(t, jwt.SigningMethodHS256.Alg(), token.Method.Alg())
}
