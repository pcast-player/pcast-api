package model

import (
	"time"
)

// Feed is a model for a podcast feed
// @model Feed
type Feed struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	URL       string    `json:"url"`
}

// CreateFeedRequest is a model for a request to create a new feed
// @model CreateFeedRequest
type CreateFeedRequest struct {
	URL string `json:"url" validate:"required,url"`
}
