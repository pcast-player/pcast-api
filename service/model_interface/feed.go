package model_interface

import (
	"context"
	"github.com/google/uuid"
	"pcast-api/store/feed"
)

type Feed interface {
	FindAll(ctx context.Context) ([]feed.Feed, error)
	Create(ctx context.Context, feed *feed.Feed) error
	FindByID(ctx context.Context, id uuid.UUID) (*feed.Feed, error)
	Delete(ctx context.Context, feed *feed.Feed) error
	Update(ctx context.Context, feed *feed.Feed) error
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]feed.Feed, error)
	FindByIdAndUserID(ctx context.Context, id, userID uuid.UUID) (*feed.Feed, error)
}
