package user

import (
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

func (s *Service) GetUser(id uuid.UUID) (*store.User, error) {
	return s.store.FindByID(id)
}

func (s *Service) GetUsers() ([]store.User, error) {
	return s.store.FindAll()
}

func (s *Service) CreateUser(user *store.User) error {
	return s.store.Create(user)
}

func (s *Service) UpdateUser(user *store.User) error {
	return s.store.Update(user)
}

func (s *Service) DeleteUser(id uuid.UUID) error {
	user, err := s.store.FindByID(id)
	if err != nil {
		return err
	}

	return s.store.Delete(user)
}

func (s *Service) Login(email string, password string) (string, error) {
	u, err := s.store.FindByEmail(email)
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
