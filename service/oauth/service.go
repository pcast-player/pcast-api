package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"pcast-api/config"
	modelInterface "pcast-api/service/model_interface"
	store "pcast-api/store/user"
)

var (
	ErrInvalidState        = errors.New("invalid state parameter")
	ErrFailedExchange      = errors.New("failed to exchange authorization code")
	ErrFailedUserInfo      = errors.New("failed to fetch user info from Google")
	ErrUnverifiedEmail     = errors.New("email not verified by Google")
	ErrGoogleNotConfigured = errors.New("Google OAuth is not configured")
)

// GoogleUserInfo represents the user info returned by Google's userinfo API
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

type Service struct {
	userStore        modelInterface.User
	googleConfig     *oauth2.Config
	jwtSecret        string
	jwtExpirationMin int
}

func NewService(cfg *config.Config, userStore modelInterface.User) *Service {
	var googleConfig *oauth2.Config
	if cfg.Auth.GoogleClientID != "" && cfg.Auth.GoogleClientSecret != "" {
		googleConfig = &oauth2.Config{
			ClientID:     cfg.Auth.GoogleClientID,
			ClientSecret: cfg.Auth.GoogleClientSecret,
			RedirectURL:  cfg.Auth.GoogleRedirectURL,
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     google.Endpoint,
		}
	}

	return &Service{
		userStore:        userStore,
		googleConfig:     googleConfig,
		jwtSecret:        cfg.Auth.JwtSecret,
		jwtExpirationMin: cfg.Auth.JwtExpirationMin,
	}
}

// GetGoogleAuthURL returns the URL to redirect user to Google's OAuth consent screen
func (s *Service) GetGoogleAuthURL(state string) (string, error) {
	if s.googleConfig == nil {
		return "", ErrGoogleNotConfigured
	}
	return s.googleConfig.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

// HandleGoogleCallback exchanges the authorization code for tokens,
// fetches user info, and creates or links the user account
func (s *Service) HandleGoogleCallback(ctx context.Context, code, state string) (string, error) {
	if s.googleConfig == nil {
		return "", ErrGoogleNotConfigured
	}

	// Exchange authorization code for access token
	token, err := s.googleConfig.Exchange(ctx, code)
	if err != nil {
		return "", ErrFailedExchange
	}

	// Fetch user info from Google
	userInfo, err := s.fetchGoogleUserInfo(ctx, token)
	if err != nil {
		return "", err
	}

	// Ensure the email is verified
	if !userInfo.VerifiedEmail {
		return "", ErrUnverifiedEmail
	}

	// Try to find existing user by Google ID
	user, err := s.userStore.FindByGoogleID(ctx, userInfo.ID)
	if err == nil && user != nil {
		// User exists with this Google ID, create JWT and return
		return s.createJwtToken(user)
	}

	// Try to find existing user by email (for account linking)
	user, err = s.userStore.FindByEmail(ctx, userInfo.Email)
	if err == nil && user != nil {
		// User exists with this email, link Google account
		if err := s.userStore.UpdateGoogleID(ctx, user.ID, userInfo.ID); err != nil {
			return "", err
		}
		return s.createJwtToken(user)
	}

	// Create new OAuth user
	newUser := &store.User{
		Email:    userInfo.Email,
		GoogleID: &userInfo.ID,
	}

	if err := s.userStore.CreateOAuthUser(ctx, newUser); err != nil {
		return "", err
	}

	return s.createJwtToken(newUser)
}

// fetchGoogleUserInfo fetches user info from Google's userinfo API
func (s *Service) fetchGoogleUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
	client := s.googleConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, ErrFailedUserInfo
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrFailedUserInfo
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrFailedUserInfo
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, ErrFailedUserInfo
	}

	return &userInfo, nil
}

// createJwtToken creates a JWT token for the given user
func (s *Service) createJwtToken(user *store.User) (string, error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.jwtExpirationMin) * time.Minute)),
		Subject:   user.ID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateState generates a random state string for CSRF protection
func GenerateState() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
