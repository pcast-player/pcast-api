package feed

import (
	"context"
	"database/sql"
	"fmt"
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

func (s *Store) FindAll(ctx context.Context) ([]Feed, error) {
	rows, err := s.queries.FindAllFeeds(ctx)
	if err != nil {
		return nil, err
	}

	// Convert sqlc models to domain models
	feeds := make([]Feed, len(rows))
	for i, row := range rows {
		feeds[i] = convertFeedRowToModel(*row)
	}
	return feeds, nil
}

func (s *Store) FindByID(ctx context.Context, id uuid.UUID) (*Feed, error) {
	row, err := s.queries.FindFeedByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return convertFeedRowToModelPtr(*row), nil
}

func (s *Store) FindByUserID(ctx context.Context, userID uuid.UUID) ([]Feed, error) {
	rows, err := s.queries.FindFeedsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Convert sqlc models to domain models
	feeds := make([]Feed, len(rows))
	for i, row := range rows {
		feeds[i] = convertFeedRowToModel(*row)
	}
	return feeds, nil
}

func (s *Store) FindByIDAndUserID(ctx context.Context, id, userID uuid.UUID) (*Feed, error) {
	row, err := s.queries.FindFeedByIDAndUserID(ctx, sqlcgen.FindFeedByIDAndUserIDParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	return convertFeedRowToModelPtr(*row), nil
}

func (s *Store) Create(ctx context.Context, feed *Feed) error {
	if err := feed.BeforeCreate(); err != nil {
		return err
	}

	_, err := s.queries.CreateFeed(ctx, sqlcgen.CreateFeedParams{
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

func (s *Store) Update(ctx context.Context, feed *Feed) error {
	if feed.URL == "" {
		return fmt.Errorf("url cannot be empty")
	}

	feed.UpdatedAt = time.Now()

	return s.queries.UpdateFeed(ctx, sqlcgen.UpdateFeedParams{
		ID:        feed.ID,
		UpdatedAt: feed.UpdatedAt,
		UserID:    feed.UserID,
		Title:     feed.Title,
		Url:       feed.URL,
		SyncedAt:  timePtrToNullTime(feed.SyncedAt),
	})
}

func (s *Store) Delete(ctx context.Context, feed *Feed) error {
	return s.queries.DeleteFeed(ctx, feed.ID)
}

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

// Helper function to convert sqlcgen.Feed to Feed
func convertFeedRowToModel(row sqlcgen.Feed) Feed {
	return Feed{
		ID:        row.ID,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		UserID:    row.UserID,
		Title:     row.Title,
		URL:       row.Url,
		SyncedAt:  nullTimeToTimePtr(row.SyncedAt),
	}
}

// Helper function to convert sqlcgen.Feed to *Feed
func convertFeedRowToModelPtr(row sqlcgen.Feed) *Feed {
	feed := convertFeedRowToModel(row)
	return &feed
}
