package episode

import (
	"github.com/google/uuid"
	"time"
)

type Episode struct {
	ID              uuid.UUID
	CreatedAt       time.Time
	UpdatedAt       time.Time
	FeedId          uuid.UUID
	FeedGUID        string
	CurrentPosition *int
	Played          bool
}

// BeforeCreate sets default values before creating an episode
// Call this explicitly in Store.Create()
func (e *Episode) BeforeCreate() error {
	if e.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		e.ID = id
	}

	if e.CreatedAt.IsZero() {
		e.CreatedAt = time.Now()
	}

	if e.UpdatedAt.IsZero() {
		e.UpdatedAt = time.Now()
	}

	// Set default played to false if not set
	// Note: This doesn't affect zero value (false), just documenting the default

	return nil
}
