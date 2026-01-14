package user

import (
	"context"
	"errors"
	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	modelInterface "pcast-api/service/model_interface"
	store "pcast-api/store/user"
	"time"
)

type Service struct {
	store modelInterface.User
}

func NewService(store modelInterface.User) *Service {
	return &Service{store: store}
}

func (s *Service) GetUser(ctx context.Context, id uuid.UUID) (*store.User, error) {
	return s.store.FindByID(ctx, id)
}

func (s *Service) GetUsers(ctx context.Context) ([]store.User, error) {
	return s.store.FindAll(ctx)
}

func (s *Service) CreateUser(ctx context.Context, user *store.User) error {
	return s.store.Create(ctx, user)
}

func (s *Service) UpdateUser(ctx context.Context, user *store.User) error {
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
		return "", err
	}

	match, err := argon2id.ComparePasswordAndHash(password, u.Password)
	if err != nil {
		return "", err
	}
	if !match {
		return "", errors.New("invalid password")
	}

	return createJwtToken(u)
}

func createJwtToken(user *store.User) (string, error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
		Subject:   user.ID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte("testsecret"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
