package feed

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"pcast-api/db/sqlcgen"
)

type Store struct {
	db      *sql.DB
	queries *sqlcgen.Queries
}

func New(database *sql.DB) *Store {
	return &Store{
		db:      database,
		queries: sqlcgen.New(database),
	}
}

func (s *Store) FindAll() ([]Feed, error) {
	rows, err := s.queries.FindAllFeeds(context.Background())
	if err != nil {
		return nil, err
	}

	// Convert sqlc models to domain models
	feeds := make([]Feed, len(rows))
	for i, row := range rows {
		feeds[i] = Feed{
			ID:        row.ID,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
			UserID:    row.UserID,
			Title:     row.Title,
			URL:       row.Url,
			SyncedAt:  nullTimeToTimePtr(row.SyncedAt),
		}
	}
	return feeds, nil
}

func (s *Store) FindByID(id uuid.UUID) (*Feed, error) {
	row, err := s.queries.FindFeedByID(context.Background(), id)
	if err != nil {
		return nil, err
	}

	return &Feed{
		ID:        row.ID,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		UserID:    row.UserID,
		Title:     row.Title,
		URL:       row.Url,
		SyncedAt:  nullTimeToTimePtr(row.SyncedAt),
	}, nil
}

func (s *Store) FindByUserID(userID uuid.UUID) ([]Feed, error) {
	rows, err := s.queries.FindFeedsByUserID(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	// Convert sqlc models to domain models
	feeds := make([]Feed, len(rows))
	for i, row := range rows {
		feeds[i] = Feed{
			ID:        row.ID,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
			UserID:    row.UserID,
			Title:     row.Title,
			URL:       row.Url,
			SyncedAt:  nullTimeToTimePtr(row.SyncedAt),
		}
	}
	return feeds, nil
}

func (s *Store) FindByIdAndUserID(id, userID uuid.UUID) (*Feed, error) {
	row, err := s.queries.FindFeedByIDAndUserID(context.Background(), sqlcgen.FindFeedByIDAndUserIDParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	return &Feed{
		ID:        row.ID,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		UserID:    row.UserID,
		Title:     row.Title,
		URL:       row.Url,
		SyncedAt:  nullTimeToTimePtr(row.SyncedAt),
	}, nil
}

func (s *Store) Create(feed *Feed) error {
	if err := feed.BeforeCreate(); err != nil {
		return err
	}

	_, err := s.queries.CreateFeed(context.Background(), sqlcgen.CreateFeedParams{
		ID:        feed.ID,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
		UserID:    feed.UserID,
		Title:     feed.Title,
		Url:       feed.URL,
		SyncedAt:  timePtrToNullTime(feed.SyncedAt),
	})

	return err
}

func (s *Store) Update(feed *Feed) error {
	feed.UpdatedAt = time.Now()

	return s.queries.UpdateFeed(context.Background(), sqlcgen.UpdateFeedParams{
		ID:        feed.ID,
		UpdatedAt: feed.UpdatedAt,
		UserID:    feed.UserID,
		Title:     feed.Title,
		Url:       feed.URL,
		SyncedAt:  timePtrToNullTime(feed.SyncedAt),
	})
}

func (s *Store) Delete(feed *Feed) error {
	return s.queries.DeleteFeed(context.Background(), feed.ID)
}

// Helper functions to convert between *time.Time and sql.NullTime
func timePtrToNullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: *t, Valid: true}
}

func nullTimeToTimePtr(nt sql.NullTime) *time.Time {
	if !nt.Valid {
		return nil
	}
	return &nt.Time
}
