package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"

	store "pcast-api/store/user"
)

// Helper function to create a pointer to a string
func strPtr(s string) *string {
	return &s
}

// mockUserStore implements modelInterface.User for testing
type mockUserStore struct {
	user              *store.User
	users             []store.User
	findByGoogleIDErr error
	findByEmailErr    error
	createErr         error
	updateErr         error
}

func (m *mockUserStore) FindAll(ctx context.Context) ([]store.User, error) {
	return m.users, nil
}

func (m *mockUserStore) FindByID(ctx context.Context, id uuid.UUID) (*store.User, error) {
	return m.user, nil
}

func (m *mockUserStore) FindByEmail(ctx context.Context, email string) (*store.User, error) {
	if m.findByEmailErr != nil {
		return nil, m.findByEmailErr
	}
	return m.user, nil
}

func (m *mockUserStore) FindByGoogleID(ctx context.Context, googleID string) (*store.User, error) {
	if m.findByGoogleIDErr != nil {
		return nil, m.findByGoogleIDErr
	}
	return m.user, nil
}

func (m *mockUserStore) Create(ctx context.Context, user *store.User) error {
	return m.createErr
}

func (m *mockUserStore) CreateOAuthUser(ctx context.Context, user *store.User) error {
	if m.createErr != nil {
		return m.createErr
	}
	// Simulate BeforeCreate setting the ID
	if user.ID == uuid.Nil {
		user.ID = uuid.Must(uuid.NewV7())
	}
	return nil
}

func (m *mockUserStore) Update(ctx context.Context, user *store.User) error {
	return m.updateErr
}

func (m *mockUserStore) UpdateGoogleID(ctx context.Context, userID uuid.UUID, googleID string) error {
	return m.updateErr
}

func (m *mockUserStore) Delete(ctx context.Context, user *store.User) error {
	return nil
}

// mockOAuthProvider implements OAuthProvider for testing
type mockOAuthProvider struct {
	authCodeURL  string
	token        *oauth2.Token
	exchangeErr  error
	httpClient   *http.Client
}

func (m *mockOAuthProvider) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
	return m.authCodeURL
}

func (m *mockOAuthProvider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	if m.exchangeErr != nil {
		return nil, m.exchangeErr
	}
	return m.token, nil
}

func (m *mockOAuthProvider) Client(ctx context.Context, token *oauth2.Token) *http.Client {
	return m.httpClient
}

func TestGenerateState(t *testing.T) {
	state, err := GenerateState()
	assert.NoError(t, err)
	assert.NotEmpty(t, state)

	// Should be a valid UUID
	_, err = uuid.Parse(state)
	assert.NoError(t, err)
}

func TestService_GetGoogleAuthURL_Success(t *testing.T) {
	provider := &mockOAuthProvider{
		authCodeURL: "https://accounts.google.com/o/oauth2/auth?client_id=test",
	}
	service := NewServiceWithProvider(nil, provider, "secret", 60)

	url, err := service.GetGoogleAuthURL("test-state")
	assert.NoError(t, err)
	assert.Equal(t, "https://accounts.google.com/o/oauth2/auth?client_id=test", url)
}

func TestService_GetGoogleAuthURL_NotConfigured(t *testing.T) {
	service := NewServiceWithProvider(nil, nil, "secret", 60)

	url, err := service.GetGoogleAuthURL("test-state")
	assert.Error(t, err)
	assert.Equal(t, ErrGoogleNotConfigured, err)
	assert.Empty(t, url)
}

func TestService_HandleGoogleCallback_NotConfigured(t *testing.T) {
	service := NewServiceWithProvider(nil, nil, "secret", 60)

	token, err := service.HandleGoogleCallback(context.Background(), "code")
	assert.Error(t, err)
	assert.Equal(t, ErrGoogleNotConfigured, err)
	assert.Empty(t, token)
}

func TestService_HandleGoogleCallback_ExchangeFails(t *testing.T) {
	provider := &mockOAuthProvider{
		exchangeErr: errors.New("exchange failed"),
	}
	service := NewServiceWithProvider(nil, provider, "secret", 60)

	token, err := service.HandleGoogleCallback(context.Background(), "invalid-code")
	assert.Error(t, err)
	assert.Equal(t, ErrFailedExchange, err)
	assert.Empty(t, token)
}

func TestService_HandleGoogleCallback_ExistingGoogleUser(t *testing.T) {
	userID := uuid.Must(uuid.NewV7())
	existingUser := &store.User{
		ID:       userID,
		Email:    "test@example.com",
		GoogleID: strPtr("google123"),
	}

	userStore := &mockUserStore{
		user: existingUser,
	}

	// Create a test server that returns Google user info
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userInfo := GoogleUserInfo{
			ID:            "google123",
			Email:         "test@example.com",
			VerifiedEmail: true,
			Name:          "Test User",
		}
		json.NewEncoder(w).Encode(userInfo)
	}))
	defer server.Close()

	provider := &mockOAuthProvider{
		token:      &oauth2.Token{AccessToken: "test-token"},
		httpClient: server.Client(),
	}
	// Override the client to use our test server
	provider.httpClient = &http.Client{
		Transport: &testTransport{server: server},
	}

	service := NewServiceWithProvider(userStore, provider, "secret", 60)

	token, err := service.HandleGoogleCallback(context.Background(), "valid-code")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestService_HandleGoogleCallback_LinkExistingEmailUser(t *testing.T) {
	userID := uuid.Must(uuid.NewV7())
	existingUser := &store.User{
		ID:       userID,
		Email:    "test@example.com",
		Password: strPtr("hashedpassword"),
	}

	userStore := &mockUserStore{
		user:              existingUser,
		findByGoogleIDErr: errors.New("not found"), // No Google ID yet
	}

	// Create a test server that returns Google user info
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userInfo := GoogleUserInfo{
			ID:            "google456",
			Email:         "test@example.com",
			VerifiedEmail: true,
			Name:          "Test User",
		}
		json.NewEncoder(w).Encode(userInfo)
	}))
	defer server.Close()

	provider := &mockOAuthProvider{
		token: &oauth2.Token{AccessToken: "test-token"},
		httpClient: &http.Client{
			Transport: &testTransport{server: server},
		},
	}

	service := NewServiceWithProvider(userStore, provider, "secret", 60)

	token, err := service.HandleGoogleCallback(context.Background(), "valid-code")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestService_HandleGoogleCallback_NewUser(t *testing.T) {
	userStore := &mockUserStore{
		findByGoogleIDErr: errors.New("not found"),
		findByEmailErr:    errors.New("not found"),
	}

	// Create a test server that returns Google user info
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userInfo := GoogleUserInfo{
			ID:            "google789",
			Email:         "newuser@example.com",
			VerifiedEmail: true,
			Name:          "New User",
		}
		json.NewEncoder(w).Encode(userInfo)
	}))
	defer server.Close()

	provider := &mockOAuthProvider{
		token: &oauth2.Token{AccessToken: "test-token"},
		httpClient: &http.Client{
			Transport: &testTransport{server: server},
		},
	}

	service := NewServiceWithProvider(userStore, provider, "secret", 60)

	token, err := service.HandleGoogleCallback(context.Background(), "valid-code")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestService_HandleGoogleCallback_UnverifiedEmail(t *testing.T) {
	userStore := &mockUserStore{}

	// Create a test server that returns unverified Google user info
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userInfo := GoogleUserInfo{
			ID:            "google123",
			Email:         "unverified@example.com",
			VerifiedEmail: false,
			Name:          "Unverified User",
		}
		json.NewEncoder(w).Encode(userInfo)
	}))
	defer server.Close()

	provider := &mockOAuthProvider{
		token: &oauth2.Token{AccessToken: "test-token"},
		httpClient: &http.Client{
			Transport: &testTransport{server: server},
		},
	}

	service := NewServiceWithProvider(userStore, provider, "secret", 60)

	token, err := service.HandleGoogleCallback(context.Background(), "valid-code")
	assert.Error(t, err)
	assert.Equal(t, ErrUnverifiedEmail, err)
	assert.Empty(t, token)
}

func TestService_HandleGoogleCallback_CreateUserFails(t *testing.T) {
	userStore := &mockUserStore{
		findByGoogleIDErr: errors.New("not found"),
		findByEmailErr:    errors.New("not found"),
		createErr:         errors.New("database error"),
	}

	// Create a test server that returns Google user info
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userInfo := GoogleUserInfo{
			ID:            "google123",
			Email:         "test@example.com",
			VerifiedEmail: true,
			Name:          "Test User",
		}
		json.NewEncoder(w).Encode(userInfo)
	}))
	defer server.Close()

	provider := &mockOAuthProvider{
		token: &oauth2.Token{AccessToken: "test-token"},
		httpClient: &http.Client{
			Transport: &testTransport{server: server},
		},
	}

	service := NewServiceWithProvider(userStore, provider, "secret", 60)

	token, err := service.HandleGoogleCallback(context.Background(), "valid-code")
	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
	assert.Empty(t, token)
}

func TestService_HandleGoogleCallback_GoogleAPIFails(t *testing.T) {
	userStore := &mockUserStore{}

	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	provider := &mockOAuthProvider{
		token: &oauth2.Token{AccessToken: "test-token"},
		httpClient: &http.Client{
			Transport: &testTransport{server: server},
		},
	}

	service := NewServiceWithProvider(userStore, provider, "secret", 60)

	token, err := service.HandleGoogleCallback(context.Background(), "valid-code")
	assert.Error(t, err)
	assert.Equal(t, ErrFailedUserInfo, err)
	assert.Empty(t, token)
}

// testTransport redirects all requests to the test server
type testTransport struct {
	server *httptest.Server
}

func (t *testTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Redirect the request to our test server
	req.URL.Scheme = "http"
	req.URL.Host = t.server.Listener.Addr().String()
	return http.DefaultTransport.RoundTrip(req)
}
