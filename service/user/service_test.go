package user

import (
	"context"
	"fmt"
	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	store "pcast-api/store/user"
	"testing"
)

type mockStore struct {
	user *store.User
	err  error
}

func (m *mockStore) FindByEmail(ctx context.Context, email string) (*store.User, error) {
	return m.user, m.err
}

func (m *mockStore) FindByID(ctx context.Context, id uuid.UUID) (*store.User, error) {
	return m.user, m.err
}

func (m *mockStore) FindAll(ctx context.Context) ([]store.User, error) {
	return []store.User{*m.user}, m.err
}

func (m *mockStore) Create(ctx context.Context, user *store.User) error {
	return m.err
}

func (m *mockStore) Update(ctx context.Context, user *store.User) error {
	return m.err
}

func (m *mockStore) Delete(ctx context.Context, user *store.User) error {
	return m.err
}

func TestService_GetUser(t *testing.T) {
	user := &store.User{Email: "foo@bar.com", Password: "password"}
	s := &mockStore{user: user}
	service := NewService(s, "testsecret", 10)

	result, err := service.GetUser(context.Background(), user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user, result)
}

func TestService_GetUsers(t *testing.T) {
	user := &store.User{Email: "foo@bar.com", Password: "password"}
	s := &mockStore{user: user}
	service := NewService(s, "testsecret", 10)

	result, err := service.GetUsers(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, []store.User{*user}, result)
}

func TestService_CreateUser(t *testing.T) {
	user := &store.User{Email: "foo@bar.com", Password: "password"}
	s := &mockStore{user: user}
	service := NewService(s, "testsecret", 10)

	result, err := service.CreateUser(context.Background(), user.Email, "password")
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestService_UpdateUser(t *testing.T) {
	user := &store.User{Email: "foo@bar.com", Password: "password"}
	s := &mockStore{user: user}
	service := NewService(s, "testsecret", 10)

	err := service.UpdateUser(context.Background(), user)
	assert.NoError(t, err)
}

func TestService_DeleteUser(t *testing.T) {
	user := &store.User{Email: "foo@bar.com", Password: "password"}
	s := &mockStore{user: user}
	service := NewService(s, "testsecret", 10)

	err := service.DeleteUser(context.Background(), user.ID)
	assert.NoError(t, err)
}

func TestService_DeleteUser_Error(t *testing.T) {
	user := &store.User{Email: "foo@bar.com", Password: "password"}
	s := &mockStore{user: user, err: assert.AnError}
	service := NewService(s, "testsecret", 10)

	err := service.DeleteUser(context.Background(), user.ID)
	assert.Error(t, err)
}

func TestService_CreateUser_Error(t *testing.T) {
	user := &store.User{Email: "foo@bar.com", Password: "password"}
	s := &mockStore{user: user, err: assert.AnError}
	service := NewService(s, "testsecret", 10)

	_, err := service.CreateUser(context.Background(), user.Email, "password")
	assert.Error(t, err)
}

func TestService_UpdateUser_Error(t *testing.T) {
	user := &store.User{Email: "foo@bar.com", Password: "password"}
	s := &mockStore{user: user, err: assert.AnError}
	service := NewService(s, "testsecret", 10)

	err := service.UpdateUser(context.Background(), user)
	assert.Error(t, err)
}

func TestService_GetUser_Error(t *testing.T) {
	user := &store.User{Email: "foo@bar.com", Password: "password"}
	s := &mockStore{user: user, err: assert.AnError}
	service := NewService(s, "testsecret", 10)

	_, err := service.GetUser(context.Background(), user.ID)
	assert.Error(t, err)
}

func TestService_GetUsers_Error(t *testing.T) {
	user := &store.User{Email: "foo@bar.com", Password: "password"}
	s := &mockStore{user: user, err: assert.AnError}
	service := NewService(s, "testsecret", 10)

	_, err := service.GetUsers(context.Background())
	assert.Error(t, err)
}

func TestService_Login(t *testing.T) {
	userID, err := uuid.NewV7()
	assert.NoError(t, err)

	password := "password"
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	assert.NoError(t, err)

	user := &store.User{ID: userID, Email: "foo@bar.com", Password: hash}
	s := &mockStore{user: user}
	service := NewService(s, "testsecret", 10)

	tokenString, err := service.Login(context.Background(), user.Email, password)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("testsecret"), nil
	})
	assert.NoError(t, err)

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		assert.Equal(t, user.ID.String(), claims["sub"])
	} else {
		assert.Fail(t, "claims not found")
	}
}

func TestService_Login_UserNotFound(t *testing.T) {
	s := &mockStore{err: assert.AnError}
	service := NewService(s, "testsecret", 10)

	tokenString, err := service.Login(context.Background(), "nonexistent@bar.com", "password")
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidPassword, err)
	assert.Empty(t, tokenString)
}

func TestService_Login_WrongPassword(t *testing.T) {
	userID := uuid.Must(uuid.NewV7())

	hash, err := argon2id.CreateHash("correctpassword", argon2id.DefaultParams)
	assert.NoError(t, err)

	user := &store.User{ID: userID, Email: "foo@bar.com", Password: hash}
	s := &mockStore{user: user}
	service := NewService(s, "testsecret", 10)

	tokenString, err := service.Login(context.Background(), user.Email, "wrongpassword")
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidPassword, err)
	assert.Empty(t, tokenString)
}

func TestService_UpdatePassword(t *testing.T) {
	userID := uuid.Must(uuid.NewV7())
	oldPassword := "oldpassword"
	newPassword := "newpassword"

	hash, err := argon2id.CreateHash(oldPassword, argon2id.DefaultParams)
	assert.NoError(t, err)

	user := &store.User{ID: userID, Email: "foo@bar.com", Password: hash}
	s := &mockStore{user: user}
	service := NewService(s, "testsecret", 10)

	err = service.UpdatePassword(context.Background(), userID, oldPassword, newPassword)
	assert.NoError(t, err)

	// Verify the password was updated
	match, err := argon2id.ComparePasswordAndHash(newPassword, user.Password)
	assert.NoError(t, err)
	assert.True(t, match)
}

func TestService_UpdatePassword_WrongOldPassword(t *testing.T) {
	userID := uuid.Must(uuid.NewV7())

	hash, err := argon2id.CreateHash("correctpassword", argon2id.DefaultParams)
	assert.NoError(t, err)

	user := &store.User{ID: userID, Email: "foo@bar.com", Password: hash}
	s := &mockStore{user: user}
	service := NewService(s, "testsecret", 10)

	err = service.UpdatePassword(context.Background(), userID, "wrongpassword", "newpassword")
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidPassword, err)
}

func TestService_UpdatePassword_UserNotFound(t *testing.T) {
	s := &mockStore{err: assert.AnError}
	service := NewService(s, "testsecret", 10)

	err := service.UpdatePassword(context.Background(), uuid.Must(uuid.NewV7()), "old", "new")
	assert.Error(t, err)
	assert.Equal(t, ErrUserNotFound, err)
}
