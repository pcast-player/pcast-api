package model_interface

import (
	"github.com/google/uuid"
	"pcast-api/store/feed"
)

type Feed interface {
	FindAll() ([]feed.Feed, error)
	Create(feed *feed.Feed) error
	FindByID(id uuid.UUID) (*feed.Feed, error)
	Delete(feed *feed.Feed) error
	Update(feed *feed.Feed) error
	FindByUserID(userID uuid.UUID) ([]feed.Feed, error)
	FindByIdAndUserID(id, userID uuid.UUID) (*feed.Feed, error)
}
