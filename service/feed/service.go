package feed

import (
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

func (s *Service) GetFeed(id uuid.UUID) (*store.Feed, error) {
	return s.store.FindByID(id)
}

func (s *Service) GetFeeds(userID uuid.UUID) ([]store.Feed, error) {
	return s.store.FindByUserID(userID)
}

func (s *Service) GetFeedsByUserID(userID uuid.UUID) ([]store.Feed, error) {
	return s.store.FindByUserID(userID)
}

func (s *Service) CreateFeed(feed *store.Feed) error {
	return s.store.Create(feed)
}

func (s *Service) DeleteFeed(userID uuid.UUID, id uuid.UUID) error {
	feed, err := s.store.FindByIdAndUserID(id, userID)
	if err != nil {
		return err
	}

	return s.store.Delete(feed)
}

func (s *Service) SyncFeed(userID uuid.UUID, id uuid.UUID) error {
	feed, err := s.store.FindByIdAndUserID(id, userID)
	if err != nil {
		return err
	}

	now := time.Now()
	feed.SyncedAt = &now

	return s.store.Update(feed)
}
