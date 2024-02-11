package feed

import (
	"github.com/google/uuid"
)

type Interface interface {
	FindAll() ([]Feed, error)
	Create(feed *Feed) error
	FindByID(id uuid.UUID) (*Feed, error)
	Delete(feed *Feed) error
	Update(feed *Feed) error
	TruncateTables()
	RemoveTables()
}
