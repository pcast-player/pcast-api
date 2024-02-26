package episode

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Store {
	autoMigrateEpisodeModel(db)

	return &Store{db: db}
}

func autoMigrateEpisodeModel(db *gorm.DB) {
	err := db.AutoMigrate(&Episode{})
	if err != nil {
		panic("Failed to migrate database!")
	}
}

func (s *Store) FindAll() ([]Episode, error) {
	var episodes []Episode
	if err := s.db.Find(&episodes).Error; err != nil {
		return nil, err
	}

	return episodes, nil
}

func (s *Store) FindByID(id uuid.UUID) (*Episode, error) {
	var episode Episode
	if err := s.db.First(&episode, id).Error; err != nil {
		return nil, err
	}

	return &episode, nil
}

func (s *Store) Create(episode *Episode) error {
	return s.db.Create(episode).Error
}

func (s *Store) Update(episode *Episode) error {
	return s.db.Save(episode).Error
}

func (s *Store) Delete(episode *Episode) error {
	return s.db.Delete(episode).Error
}
