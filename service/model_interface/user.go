package model_interface

import (
	"context"
	"github.com/google/uuid"
	"pcast-api/store/user"
)

type User interface {
	FindAll(ctx context.Context) ([]user.User, error)
	Create(ctx context.Context, user *user.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*user.User, error)
	FindByEmail(ctx context.Context, email string) (*user.User, error)
	Delete(ctx context.Context, user *user.User) error
	Update(ctx context.Context, user *user.User) error
}
