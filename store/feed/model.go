package feed

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Feed struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID `gorm:"type:uuid"`
	Title     string
	URL       string
	SyncedAt  *time.Time
}

func (f *Feed) BeforeCreate(_ *gorm.DB) (err error) {
	f.ID, err = uuid.NewV7()
	if err != nil {
		return err
	}

	return nil
}
