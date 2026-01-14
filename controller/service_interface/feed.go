package service_interface

import (
	"context"
	"github.com/google/uuid"
	store "pcast-api/store/feed"
)

type Feed interface {
	GetFeed(ctx context.Context, id uuid.UUID) (*store.Feed, error)
	GetFeeds(ctx context.Context, userID uuid.UUID) ([]store.Feed, error)
	CreateFeed(ctx context.Context, feed *store.Feed) error
	DeleteFeed(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
	SyncFeed(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
	GetFeedsByUserID(ctx context.Context, userID uuid.UUID) ([]store.Feed, error)
}
