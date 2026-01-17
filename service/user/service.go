package user

import (
	"context"
	"errors"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	modelInterface "pcast-api/service/model_interface"
	store "pcast-api/store/user"
)

var (
	ErrInvalidPassword = errors.New("invalid password")
	ErrUserNotFound    = errors.New("user not found")
	ErrNoPassword      = errors.New("user has no password (OAuth account)")
)

type Service struct {
	store            modelInterface.User
	jwtSecret        string
	jwtExpirationMin int
}

func NewService(store modelInterface.User, jwtSecret string, jwtExpirationMin int) *Service {
	return &Service{store: store, jwtSecret: jwtSecret, jwtExpirationMin: jwtExpirationMin}
}

func (s *Service) GetUser(ctx context.Context, id uuid.UUID) (*store.User, error) {
	u, err := s.store.FindByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return u, nil
}

func (s *Service) GetUsers(ctx context.Context) ([]store.User, error) {
	return s.store.FindAll(ctx)
}

func (s *Service) CreateUser(ctx context.Context, email, password string) (*store.User, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return nil, err
	}

	user := &store.User{
		Email:    email,
		Password: &hash,
	}

	err = s.store.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) UpdateUser(ctx context.Context, user *store.User) error {
	return s.store.Update(ctx, user)
}

func (s *Service) UpdatePassword(ctx context.Context, userID uuid.UUID, oldPassword string, newPassword string) error {
	user, err := s.GetUser(ctx, userID)
	if err != nil {
		return err
	}

	// Check if user has a password (OAuth-only users don't)
	if user.Password == nil {
		return ErrNoPassword
	}

	match, err := argon2id.ComparePasswordAndHash(oldPassword, *user.Password)
	if err != nil {
		return err
	}
	if !match {
		return ErrInvalidPassword
	}

	hash, err := argon2id.CreateHash(newPassword, argon2id.DefaultParams)
	if err != nil {
		return err
	}

	user.Password = &hash
	return s.store.Update(ctx, user)
}

func (s *Service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	user, err := s.store.FindByID(ctx, id)
	if err != nil {
		return err
	}

	return s.store.Delete(ctx, user)
}

func (s *Service) Login(ctx context.Context, email string, password string) (string, error) {
	u, err := s.store.FindByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidPassword // Return generic error for security
	}

	// Check if user has a password (OAuth-only users can't login with password)
	if u.Password == nil {
		return "", ErrInvalidPassword
	}

	match, err := argon2id.ComparePasswordAndHash(password, *u.Password)
	if err != nil {
		return "", err
	}
	if !match {
		return "", ErrInvalidPassword
	}

	return s.createJwtToken(u)
}

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
