package model_interface

import (
	"github.com/google/uuid"
	"pcast-api/store/episode"
)

type Episode interface {
	FindAll() ([]episode.Episode, error)
	Create(episode *episode.Episode) error
	FindByID(id uuid.UUID) (*episode.Episode, error)
	Delete(episode *episode.Episode) error
	Update(episode *episode.Episode) error
}
