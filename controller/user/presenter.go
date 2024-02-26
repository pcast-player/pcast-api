package user

import (
	"github.com/google/uuid"
	store "pcast-api/store/user"
)

// Presenter represents a user presenter
// @model Presenter
type Presenter struct {
	ID uuid.UUID `json:"id"`
}

func NewPresenter(user *store.User) *Presenter {
	return &Presenter{
		ID: user.ID,
	}
}
