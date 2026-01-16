package episode

import (
	"github.com/google/uuid"
	"pcast-api/store"
	"time"
)

type Episode struct {
	ID              uuid.UUID
	CreatedAt       time.Time
	UpdatedAt       time.Time
	FeedID          uuid.UUID
	FeedGUID        string
	CurrentPosition *int
	Played          bool
}

func (e *Episode) SetID(id uuid.UUID) {
	e.ID = id
}

func (e *Episode) GetID() uuid.UUID {
	return e.ID
}

func (e *Episode) SetCreatedAt(t time.Time) {
	e.CreatedAt = t
}

func (e *Episode) GetCreatedAt() time.Time {
	return e.CreatedAt
}

func (e *Episode) SetUpdatedAt(t time.Time) {
	e.UpdatedAt = t
}

func (e *Episode) GetUpdatedAt() time.Time {
	return e.UpdatedAt
}

func (e *Episode) BeforeCreate() error {
	return store.BeforeCreate(e)
}
