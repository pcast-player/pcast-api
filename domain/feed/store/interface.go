package store

import (
	"github.com/google/uuid"
	"pcast-api/domain/feed/model"
)

type Interface interface {
	FindAll() ([]model.Feed, error)
	Create(feed *model.Feed) error
	FindByID(id uuid.UUID) (*model.Feed, error)
	Delete(feed *model.Feed) error
	Update(feed *model.Feed) error
}
