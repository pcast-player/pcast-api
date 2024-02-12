package episode

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Episode struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	FeedId          uuid.UUID `gorm:"type:uuid"`
	FeedGUID        string
	CurrentPosition *int
	Played          bool
}

func (e *Episode) BeforeCreate(_ *gorm.DB) (err error) {
	e.ID, err = uuid.NewV7()
	if err != nil {
		return err
	}

	e.Played = false

	return nil
}
