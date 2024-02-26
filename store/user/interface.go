package user

import "github.com/google/uuid"

type Interface interface {
	FindAll() ([]User, error)
	Create(user *User) error
	FindByID(id uuid.UUID) (*User, error)
	Delete(user *User) error
	Update(user *User) error
	TruncateTables()
	RemoveTables()
}
