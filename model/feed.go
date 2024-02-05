package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// Feed is a model for a podcast feed
// @model Feed
type Feed struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	URL       string    `json:"url"`
}

func (feed *Feed) BeforeCreate(tx *gorm.DB) (err error) {
	feed.ID = uuid.New()

	return nil
}

// CreateFeedRequest is a model for a request to create a new feed
// @model CreateFeedRequest
type CreateFeedRequest struct {
	URL string `json:"url" validate:"required,url"`
}
