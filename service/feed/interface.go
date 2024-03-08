package feed

import (
	"github.com/google/uuid"
	store "pcast-api/store/feed"
)

type Interface interface {
	GetFeed(id uuid.UUID) (*store.Feed, error)
	GetFeeds(userID uuid.UUID) ([]store.Feed, error)
	CreateFeed(feed *store.Feed) error
	DeleteFeed(userID uuid.UUID, id uuid.UUID) error
	SyncFeed(userID uuid.UUID, id uuid.UUID) error
	GetFeedsByUserID(userID uuid.UUID) ([]store.Feed, error)
}
