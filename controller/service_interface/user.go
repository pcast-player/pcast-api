package service_interface

import (
	"context"

	"github.com/google/uuid"

	store "pcast-api/store/user"
)

type User interface {
	CreateUser(ctx context.Context, email, password string) (*store.User, error)
	Login(ctx context.Context, email string, password string) (string, error)
	UpdatePassword(ctx context.Context, userID uuid.UUID, oldPassword string, newPassword string) error
}
