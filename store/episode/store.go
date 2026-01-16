package episode

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"pcast-api/db/sqlcgen"
)

type Store struct {
	queries *sqlcgen.Queries
}

func New(database *sql.DB) *Store {
	return &Store{
		queries: sqlcgen.New(database),
	}
}

func (s *Store) FindAll(ctx context.Context) ([]Episode, error) {
	rows, err := s.queries.FindAllEpisodes(ctx)
	if err != nil {
		return nil, err
	}

	// Convert sqlc models to domain models
	episodes := make([]Episode, len(rows))
	for i, row := range rows {
		episodes[i] = Episode{
			ID:              row.ID,
			CreatedAt:       row.CreatedAt,
			UpdatedAt:       row.UpdatedAt,
			FeedID:          row.FeedID,
			FeedGUID:        row.FeedGuid,
			CurrentPosition: nullInt32ToIntPtr(row.CurrentPosition),
			Played:          row.Played,
		}
	}
	return episodes, nil
}

func (s *Store) FindByID(ctx context.Context, id uuid.UUID) (*Episode, error) {
	row, err := s.queries.FindEpisodeByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &Episode{
		ID:              row.ID,
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
		FeedID:          row.FeedID,
		FeedGUID:        row.FeedGuid,
		CurrentPosition: nullInt32ToIntPtr(row.CurrentPosition),
		Played:          row.Played,
	}, nil
}

func (s *Store) Create(ctx context.Context, episode *Episode) error {
	if err := episode.BeforeCreate(); err != nil {
		return err
	}

	_, err := s.queries.CreateEpisode(ctx, sqlcgen.CreateEpisodeParams{
		ID:              episode.ID,
		CreatedAt:       episode.CreatedAt,
		UpdatedAt:       episode.UpdatedAt,
		FeedID:          episode.FeedID,
		FeedGuid:        episode.FeedGUID,
		CurrentPosition: intPtrToNullInt32(episode.CurrentPosition),
		Played:          episode.Played,
	})

	return err
}

func (s *Store) Update(ctx context.Context, episode *Episode) error {
	episode.UpdatedAt = time.Now()

	return s.queries.UpdateEpisode(ctx, sqlcgen.UpdateEpisodeParams{
		ID:              episode.ID,
		UpdatedAt:       episode.UpdatedAt,
		FeedID:          episode.FeedID,
		FeedGuid:        episode.FeedGUID,
		CurrentPosition: intPtrToNullInt32(episode.CurrentPosition),
		Played:          episode.Played,
	})
}

func (s *Store) Delete(ctx context.Context, episode *Episode) error {
	return s.queries.DeleteEpisode(ctx, episode.ID)
}

func intPtrToNullInt32(i *int) sql.NullInt32 {
	if i == nil {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{Int32: int32(*i), Valid: true}
}

func nullInt32ToIntPtr(n sql.NullInt32) *int {
	if !n.Valid {
		return nil
	}
	i := int(n.Int32)
	return &i
}
