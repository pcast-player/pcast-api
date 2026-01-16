package user

import (
	"github.com/google/uuid"
	"pcast-api/store"
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

func (u *User) SetID(id uuid.UUID) {
	u.ID = id
}

func (u *User) GetID() uuid.UUID {
	return u.ID
}

func (u *User) SetCreatedAt(t time.Time) {
	u.CreatedAt = t
}

func (u *User) GetCreatedAt() time.Time {
	return u.CreatedAt
}

func (u *User) SetUpdatedAt(t time.Time) {
	u.UpdatedAt = t
}

func (u *User) GetUpdatedAt() time.Time {
	return u.UpdatedAt
}

func (u *User) BeforeCreate() error {
	return store.BeforeCreate(u)
}
