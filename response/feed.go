package response

import (
	"github.com/google/uuid"
	"pcast-api/model"
	"time"
)

// Feed represents a feed response
// @model Feed
type Feed struct {
	ID       uuid.UUID  `json:"id"`
	Title    string     `json:"title"`
	URL      string     `json:"url"`
	SyncedAt *time.Time `json:"syncedAt"`
}

func NewFeed(feed *model.Feed) *Feed {
	return &Feed{
		ID:       feed.ID,
		Title:    feed.Title,
		URL:      feed.URL,
		SyncedAt: feed.SyncedAt,
	}
}
