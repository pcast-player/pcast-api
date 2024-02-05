package model

import "gorm.io/gorm"

// Feed is a model for a podcast feed
// @model Feed
type Feed struct {
	gorm.Model
	URL string
}

// CreateFeedRequest is a model for a request to create a new feed
// @model CreateFeedRequest
type CreateFeedRequest struct {
	URL string `json:"url" validate:"required,url"`
}
