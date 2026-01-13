package feed

import (
	"github.com/google/uuid"
	"time"
)

type Feed struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	Title     string
	URL       string
	SyncedAt  *time.Time
}

// BeforeCreate sets default values before creating a feed
// Call this explicitly in Store.Create()
func (f *Feed) BeforeCreate() error {
	if f.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		f.ID = id
	}

	if f.CreatedAt.IsZero() {
		f.CreatedAt = time.Now()
	}

	if f.UpdatedAt.IsZero() {
		f.UpdatedAt = time.Now()
	}

	return nil
}
