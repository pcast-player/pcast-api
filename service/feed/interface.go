package feed

import (
	"github.com/google/uuid"
	store "pcast-api/store/feed"
)

type Interface interface {
	GetFeed(id uuid.UUID) (*store.Feed, error)
	GetFeeds() ([]store.Feed, error)
	CreateFeed(feed *store.Feed) error
	DeleteFeed(id uuid.UUID) error
	SyncFeed(id uuid.UUID) error
}
