package feed

import (
	"github.com/google/uuid"
	"pcast-api/model"
)

type Store interface {
	FindAll() ([]model.Feed, error)
	Create(feed *model.Feed) error
	FindByID(id uuid.UUID) (*model.Feed, error)
	Delete(feed *model.Feed) error
}
