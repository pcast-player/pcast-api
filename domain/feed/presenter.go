package feed

import (
	"github.com/google/uuid"
	"pcast-api/store/feed"
	"time"
)

// Presenter represents a feed presenter
// @model Presenter
type Presenter struct {
	ID       uuid.UUID  `json:"id"`
	Title    string     `json:"title"`
	URL      string     `json:"url"`
	SyncedAt *time.Time `json:"syncedAt"`
}

func NewPresenter(feed *feed.Feed) *Presenter {
	return &Presenter{
		ID:       feed.ID,
		Title:    feed.Title,
		URL:      feed.URL,
		SyncedAt: feed.SyncedAt,
	}
}
