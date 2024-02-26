package feed

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Store {
	autoMigrateFeedModel(db)

	return &Store{db: db}
}

func autoMigrateFeedModel(db *gorm.DB) {
	err := db.AutoMigrate(&Feed{})
	if err != nil {
		panic("Failed to migrate database!")
	}
}

func (s *Store) FindAll() ([]Feed, error) {
	var feeds []Feed
	if err := s.db.Find(&feeds).Error; err != nil {
		return nil, err
	}

	return feeds, nil
}

func (s *Store) FindByID(id uuid.UUID) (*Feed, error) {
	var feed Feed
	if err := s.db.First(&feed, id).Error; err != nil {
		return nil, err
	}

	return &feed, nil
}

func (s *Store) Create(feed *Feed) error {
	return s.db.Create(feed).Error
}

func (s *Store) Update(feed *Feed) error {
	return s.db.Save(feed).Error
}

func (s *Store) Delete(feed *Feed) error {
	return s.db.Delete(feed).Error
}
