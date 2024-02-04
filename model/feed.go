package model

import "gorm.io/gorm"

type Feed struct {
	gorm.Model
	URL string
}

type CreateFeedRequest struct {
	URL string `json:"url" validate:"required,url"`
}
