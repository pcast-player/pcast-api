package user

import (
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

func (m *mockStore) FindByEmail(email string) (*store.User, error) {
	return m.user, m.err
}

func (m *mockStore) FindByID(id uuid.UUID) (*store.User, error) {
	return m.user, m.err
}

func (m *mockStore) FindAll() ([]store.User, error) {
	return []store.User{*m.user}, m.err
}

func (m *mockStore) Create(user *store.User) error {
	return m.err
}

func (m *mockStore) Update(user *store.User) error {
	return m.err
}

func (m *mockStore) Delete(user *store.User) error {
	return m.err
}

func TestService_GetUser(t *testing.T) {
	user := &store.User{Email: "foo@bar.com", Password: "password"}
	s := &mockStore{user: user}
	service := NewService(s)

	result, err := service.GetUser(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user, result)
}

func TestService_GetUsers(t *testing.T) {
	user := &store.User{Email: "foo@bar.com", Password: "password"}
	s := &mockStore{user: user}
	service := NewService(s)

	result, err := service.GetUsers()
	assert.NoError(t, err)
	assert.Equal(t, []store.User{*user}, result)
}

func TestService_CreateUser(t *testing.T) {
	user := &store.User{Email: "foo@bar.com", Password: "password"}
	s := &mockStore{user: user}
	service := NewService(s)

	err := service.CreateUser(user)
	assert.NoError(t, err)
}

func TestService_UpdateUser(t *testing.T) {
	user := &store.User{Email: "foo@bar.com", Password: "password"}
	s := &mockStore{user: user}
	service := NewService(s)

	err := service.UpdateUser(user)
	assert.NoError(t, err)
}

func TestService_DeleteUser(t *testing.T) {
	user := &store.User{Email: "foo@bar.com", Password: "password"}
	s := &mockStore{user: user}
	service := NewService(s)

	err := service.DeleteUser(user.ID)
	assert.NoError(t, err)
}

func TestService_DeleteUser_Error(t *testing.T) {
	user := &store.User{Email: "foo@bar.com", Password: "password"}
	s := &mockStore{user: user, err: assert.AnError}
	service := NewService(s)

	err := service.DeleteUser(user.ID)
	assert.Error(t, err)
}

func TestService_CreateUser_Error(t *testing.T) {
	user := &store.User{Email: "foo@bar.com", Password: "password"}
	s := &mockStore{user: user, err: assert.AnError}
	service := NewService(s)

	err := service.CreateUser(user)
	assert.Error(t, err)
}

func TestService_UpdateUser_Error(t *testing.T) {
	user := &store.User{Email: "foo@bar.com", Password: "password"}
	s := &mockStore{user: user, err: assert.AnError}
	service := NewService(s)

	err := service.UpdateUser(user)
	assert.Error(t, err)
}

func TestService_GetUser_Error(t *testing.T) {
	user := &store.User{Email: "foo@bar.com", Password: "password"}
	s := &mockStore{user: user, err: assert.AnError}
	service := NewService(s)

	_, err := service.GetUser(user.ID)
	assert.Error(t, err)
}

func TestService_GetUsers_Error(t *testing.T) {
	user := &store.User{Email: "foo@bar.com", Password: "password"}
	s := &mockStore{user: user, err: assert.AnError}
	service := NewService(s)

	_, err := service.GetUsers()
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
	service := NewService(s)

	tokenString, err := service.Login(user.Email, password)
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
