package service_interface

import (
	"context"

	"github.com/google/uuid"
	store "pcast-api/store/user"
)

type User interface {
	GetUser(ctx context.Context, id uuid.UUID) (*store.User, error)
	GetUsers(ctx context.Context) ([]store.User, error)
	CreateUser(ctx context.Context, user *store.User) error
	UpdateUser(ctx context.Context, user *store.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
