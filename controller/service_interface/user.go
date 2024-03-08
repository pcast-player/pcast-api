package service_interface

import (
	"github.com/google/uuid"
	store "pcast-api/store/user"
)

type User interface {
	GetUser(id uuid.UUID) (*store.User, error)
	GetUsers() ([]store.User, error)
	CreateUser(user *store.User) error
	UpdateUser(user *store.User) error
	DeleteUser(id uuid.UUID) error
}
