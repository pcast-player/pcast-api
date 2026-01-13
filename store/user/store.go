package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"pcast-api/db/sqlcgen"
)

type Store struct {
	db      *sql.DB
	queries *sqlcgen.Queries
}

func New(database *sql.DB) *Store {
	return &Store{
		db:      database,
		queries: sqlcgen.New(database),
	}
}

func (s *Store) FindAll() ([]User, error) {
	rows, err := s.queries.FindAllUsers(context.Background())
	if err != nil {
		return nil, err
	}

	// Convert sqlc models to domain models
	users := make([]User, len(rows))
	for i, row := range rows {
		users[i] = User{
			ID:        row.ID,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
			Email:     row.Email,
			Password:  row.Password,
			Feeds:     nil, // Not loaded by default
		}
	}
	return users, nil
}

func (s *Store) FindByID(id uuid.UUID) (*User, error) {
	row, err := s.queries.FindUserByID(context.Background(), id)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:        row.ID,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		Email:     row.Email,
		Password:  row.Password,
		Feeds:     nil, // Not loaded by default
	}, nil
}

func (s *Store) FindByEmail(email string) (*User, error) {
	row, err := s.queries.FindUserByEmail(context.Background(), email)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:        row.ID,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		Email:     row.Email,
		Password:  row.Password,
		Feeds:     nil, // Not loaded by default
	}, nil
}

func (s *Store) Create(user *User) error {
	if err := user.BeforeCreate(); err != nil {
		return err
	}

	_, err := s.queries.CreateUser(context.Background(), sqlcgen.CreateUserParams{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Password:  user.Password,
	})

	return err
}

func (s *Store) Update(user *User) error {
	user.UpdatedAt = time.Now()

	return s.queries.UpdateUser(context.Background(), sqlcgen.UpdateUserParams{
		ID:        user.ID,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Password:  user.Password,
	})
}

func (s *Store) Delete(user *User) error {
	return s.queries.DeleteUser(context.Background(), user.ID)
}
