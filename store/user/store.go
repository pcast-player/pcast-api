package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"pcast-api/db/sqlcgen"
)

type Store struct {
	queries *sqlcgen.Queries
}

func New(database *sql.DB) *Store {
	return &Store{
		queries: sqlcgen.New(database),
	}
}

func (s *Store) FindAll(ctx context.Context) ([]User, error) {
	rows, err := s.queries.FindAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	// Convert sqlc models to domain models
	users := make([]User, len(rows))
	for i, row := range rows {
		users[i] = convertUserRowToModel(*row)
	}
	return users, nil
}

func (s *Store) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	row, err := s.queries.FindUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return convertUserRowToModelPtr(*row), nil
}

func (s *Store) FindByEmail(ctx context.Context, email string) (*User, error) {
	row, err := s.queries.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return convertUserRowToModelPtr(*row), nil
}

func (s *Store) Create(ctx context.Context, user *User) error {
	if err := user.BeforeCreate(); err != nil {
		return err
	}

	_, err := s.queries.CreateUser(ctx, sqlcgen.CreateUserParams{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Password:  user.Password,
	})

	return err
}

func (s *Store) Update(ctx context.Context, user *User) error {
	user.UpdatedAt = time.Now()

	return s.queries.UpdateUser(ctx, sqlcgen.UpdateUserParams{
		ID:        user.ID,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Password:  user.Password,
	})
}

func (s *Store) Delete(ctx context.Context, user *User) error {
	return s.queries.DeleteUser(ctx, user.ID)
}

// Helper function to convert sqlcgen.User to User
func convertUserRowToModel(row sqlcgen.User) User {
	return User{
		ID:        row.ID,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		Email:     row.Email,
		Password:  row.Password,
	}
}

// Helper function to convert sqlcgen.User to *User
func convertUserRowToModelPtr(row sqlcgen.User) *User {
	user := convertUserRowToModel(row)
	return &user
}
