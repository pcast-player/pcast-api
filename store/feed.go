package store

import (
	"gorm.io/gorm"
	"pcast-api/model"
)

type FeedStore struct {
	db *gorm.DB
}

func New(db *gorm.DB) *FeedStore {
	return &FeedStore{db: db}
}

func (s *FeedStore) FindAll() ([]model.Feed, error) {
	var feeds []model.Feed

	if err := s.db.Find(&feeds).Error; err != nil {
		return nil, err
	}

	return feeds, nil
}

func (s *FeedStore) Create(feed *model.Feed) error {
	return s.db.Create(feed).Error
}
