package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"pcast-api/config"
	"pcast-api/service/auth"
	modelInterface "pcast-api/service/model_interface"
	store "pcast-api/store/user"
)

var (
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

// OAuthProvider abstracts OAuth2 operations for testability
type OAuthProvider interface {
	AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
	Exchange(ctx context.Context, code string) (*oauth2.Token, error)
	Client(ctx context.Context, token *oauth2.Token) *http.Client
}

// oauth2ConfigAdapter adapts oauth2.Config to OAuthProvider interface
type oauth2ConfigAdapter struct {
	config *oauth2.Config
}

func (a *oauth2ConfigAdapter) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
	return a.config.AuthCodeURL(state, opts...)
}

func (a *oauth2ConfigAdapter) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return a.config.Exchange(ctx, code)
}

func (a *oauth2ConfigAdapter) Client(ctx context.Context, token *oauth2.Token) *http.Client {
	return a.config.Client(ctx, token)
}

type Service struct {
	userStore        modelInterface.User
	googleProvider   OAuthProvider
	jwtSecret        string
	jwtExpirationMin int
}

func NewService(cfg *config.Config, userStore modelInterface.User) *Service {
	var googleProvider OAuthProvider
	if cfg.Auth.GoogleClientID != "" && cfg.Auth.GoogleClientSecret != "" {
		googleProvider = &oauth2ConfigAdapter{
			config: &oauth2.Config{
				ClientID:     cfg.Auth.GoogleClientID,
				ClientSecret: cfg.Auth.GoogleClientSecret,
				RedirectURL:  cfg.Auth.GoogleRedirectURL,
				Scopes:       []string{"openid", "email", "profile"},
				Endpoint:     google.Endpoint,
			},
		}
	}

	return &Service{
		userStore:        userStore,
		googleProvider:   googleProvider,
		jwtSecret:        cfg.Auth.JwtSecret,
		jwtExpirationMin: cfg.Auth.JwtExpirationMin,
	}
}

// NewServiceWithProvider creates a service with a custom OAuth provider (for testing)
func NewServiceWithProvider(userStore modelInterface.User, provider OAuthProvider, jwtSecret string, jwtExpirationMin int) *Service {
	return &Service{
		userStore:        userStore,
		googleProvider:   provider,
		jwtSecret:        jwtSecret,
		jwtExpirationMin: jwtExpirationMin,
	}
}

// GetGoogleAuthURL returns the URL to redirect user to Google's OAuth consent screen
func (s *Service) GetGoogleAuthURL(state string) (string, error) {
	if s.googleProvider == nil {
		return "", ErrGoogleNotConfigured
	}
	return s.googleProvider.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

// HandleGoogleCallback exchanges the authorization code for tokens,
// fetches user info, and creates or links the user account
func (s *Service) HandleGoogleCallback(ctx context.Context, code string) (string, error) {
	if s.googleProvider == nil {
		return "", ErrGoogleNotConfigured
	}

	// Exchange authorization code for access token
	token, err := s.googleProvider.Exchange(ctx, code)
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
		return auth.CreateJWTToken(user.ID, s.jwtSecret, s.jwtExpirationMin)
	}

	// Try to find existing user by email (for account linking)
	user, err = s.userStore.FindByEmail(ctx, userInfo.Email)
	if err == nil && user != nil {
		// User exists with this email, link Google account
		if err := s.userStore.UpdateGoogleID(ctx, user.ID, userInfo.ID); err != nil {
			return "", err
		}
		return auth.CreateJWTToken(user.ID, s.jwtSecret, s.jwtExpirationMin)
	}

	// Create new OAuth user
	newUser := &store.User{
		Email:    userInfo.Email,
		GoogleID: &userInfo.ID,
	}

	if err := s.userStore.CreateOAuthUser(ctx, newUser); err != nil {
		return "", err
	}

	return auth.CreateJWTToken(newUser.ID, s.jwtSecret, s.jwtExpirationMin)
}

// fetchGoogleUserInfo fetches user info from Google's userinfo API
func (s *Service) fetchGoogleUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
	client := s.googleProvider.Client(ctx, token)
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

// GenerateState generates a random state string for CSRF protection
func GenerateState() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
