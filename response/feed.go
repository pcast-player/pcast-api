package response

import (
	"github.com/google/uuid"
	"pcast-api/model"
)

// Feed represents a feed response
// @model Feed
type Feed struct {
	ID  uuid.UUID `json:"id"`
	URL string    `json:"url"`
}

func NewFeed(feed *model.Feed) *Feed {
	return &Feed{
		ID:  feed.ID,
		URL: feed.URL,
	}
}
