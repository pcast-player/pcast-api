package feed

import (
	"github.com/google/uuid"
	store "pcast-api/store/feed"
	"time"
)

type Service struct {
	store store.Interface
}

func NewService(store store.Interface) *Service {
	return &Service{store: store}
}

func (s *Service) GetFeed(id uuid.UUID) (*store.Feed, error) {
	return s.store.FindByID(id)
}

func (s *Service) GetFeeds() ([]store.Feed, error) {
	return s.store.FindAll()
}

func (s *Service) CreateFeed(feed *store.Feed) error {
	return s.store.Create(feed)
}

func (s *Service) DeleteFeed(id uuid.UUID) error {
	feed, err := s.store.FindByID(id)
	if err != nil {
		return err
	}

	return s.store.Delete(feed)
}

func (s *Service) SyncFeed(id uuid.UUID) error {
	feed, err := s.store.FindByID(id)
	if err != nil {
		return err
	}

	now := time.Now()
	feed.SyncedAt = &now

	return s.store.Update(feed)
}
