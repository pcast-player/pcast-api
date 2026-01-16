package store

import (
	"github.com/google/uuid"
	"time"
)

type Entity interface {
	SetID(uuid.UUID)
	GetID() uuid.UUID
	SetCreatedAt(time.Time)
	GetCreatedAt() time.Time
	SetUpdatedAt(time.Time)
	GetUpdatedAt() time.Time
}

func BeforeCreate(entity Entity) error {
	if entity.GetID() == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		entity.SetID(id)
	}

	if entity.GetCreatedAt().IsZero() {
		entity.SetCreatedAt(time.Now())
	}

	if entity.GetUpdatedAt().IsZero() {
		entity.SetUpdatedAt(time.Now())
	}

	return nil
}
