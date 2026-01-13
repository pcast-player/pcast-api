package user

import (
	"github.com/google/uuid"
	"pcast-api/store/feed"
	"time"
)

type User struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Email     string
	Password  string
	Feeds     []feed.Feed `db:"-"` // Not populated by sqlc, manual load if needed
}

// BeforeCreate sets default values before creating a user
// Call this explicitly in Store.Create()
func (u *User) BeforeCreate() error {
	if u.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		u.ID = id
	}

	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}

	if u.UpdatedAt.IsZero() {
		u.UpdatedAt = time.Now()
	}

	return nil
}
