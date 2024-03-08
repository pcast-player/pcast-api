package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"pcast-api/store/feed"
	"time"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Email     string
	Password  string
	Feeds     []feed.Feed
}

func (u *User) BeforeCreate(_ *gorm.DB) (err error) {
	u.ID, err = uuid.NewV7()
	if err != nil {
		return err
	}

	return nil
}
