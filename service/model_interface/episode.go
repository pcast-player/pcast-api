package model_interface

import (
	"context"

	"github.com/google/uuid"
	"pcast-api/store/episode"
)

type Episode interface {
	FindAll(ctx context.Context) ([]episode.Episode, error)
	Create(ctx context.Context, episode *episode.Episode) error
	FindByID(ctx context.Context, id uuid.UUID) (*episode.Episode, error)
	Delete(ctx context.Context, episode *episode.Episode) error
	Update(ctx context.Context, episode *episode.Episode) error
}
