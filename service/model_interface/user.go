package model_interface

import (
	"github.com/google/uuid"
	"pcast-api/store/user"
)

type User interface {
	FindAll() ([]user.User, error)
	Create(user *user.User) error
	FindByID(id uuid.UUID) (*user.User, error)
	FindByEmail(email string) (*user.User, error)
	Delete(user *user.User) error
	Update(user *user.User) error
}
