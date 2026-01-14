package feed

import (
	"context"
	"github.com/google/uuid"
	modelInterface "pcast-api/service/model_interface"
	store "pcast-api/store/feed"
	"time"
)

type Service struct {
	store modelInterface.Feed
}

func NewService(store modelInterface.Feed) *Service {
	return &Service{store: store}
}

func (s *Service) GetFeed(ctx context.Context, id uuid.UUID) (*store.Feed, error) {
	return s.store.FindByID(ctx, id)
}

func (s *Service) GetFeeds(ctx context.Context, userID uuid.UUID) ([]store.Feed, error) {
	return s.store.FindByUserID(ctx, userID)
}

func (s *Service) GetFeedsByUserID(ctx context.Context, userID uuid.UUID) ([]store.Feed, error) {
	return s.store.FindByUserID(ctx, userID)
}

func (s *Service) CreateFeed(ctx context.Context, feed *store.Feed) error {
	return s.store.Create(ctx, feed)
}

func (s *Service) DeleteFeed(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	feed, err := s.store.FindByIdAndUserID(ctx, id, userID)
	if err != nil {
		return err
	}

	return s.store.Delete(ctx, feed)
}

func (s *Service) SyncFeed(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	feed, err := s.store.FindByIdAndUserID(ctx, id, userID)
	if err != nil {
		return err
	}

	now := time.Now()
	feed.SyncedAt = &now

	return s.store.Update(ctx, feed)
}
