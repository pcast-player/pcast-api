package feed

import (
	"github.com/google/uuid"
	"pcast-api/store"
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

func (f *Feed) SetID(id uuid.UUID) {
	f.ID = id
}

func (f *Feed) GetID() uuid.UUID {
	return f.ID
}

func (f *Feed) SetCreatedAt(t time.Time) {
	f.CreatedAt = t
}

func (f *Feed) GetCreatedAt() time.Time {
	return f.CreatedAt
}

func (f *Feed) SetUpdatedAt(t time.Time) {
	f.UpdatedAt = t
}

func (f *Feed) GetUpdatedAt() time.Time {
	return f.UpdatedAt
}

func (f *Feed) BeforeCreate() error {
	return store.BeforeCreate(f)
}
