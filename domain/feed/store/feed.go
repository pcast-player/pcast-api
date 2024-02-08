package store

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"pcast-api/domain/feed/model"
)

type Store struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s *Store) FindAll() ([]model.Feed, error) {
	var feeds []model.Feed
	if err := s.db.Find(&feeds).Error; err != nil {
		return nil, err
	}

	return feeds, nil
}

func (s *Store) Create(feed *model.Feed) error {
	return s.db.Create(feed).Error
}

func (s *Store) FindByID(id uuid.UUID) (*model.Feed, error) {
	var feed model.Feed
	if err := s.db.First(&feed, id).Error; err != nil {
		return nil, err
	}

	return &feed, nil
}

func (s *Store) Delete(feed *model.Feed) error {
	return s.db.Delete(feed).Error
}

func (s *Store) Update(feed *model.Feed) error {
	return s.db.Save(feed).Error
}
