package user

import (
	"github.com/google/uuid"
	store "pcast-api/store/user"
)

type Service struct {
	store store.Interface
}

func NewService(store store.Interface) *Service {
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
