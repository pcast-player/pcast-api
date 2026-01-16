package feed

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	store "pcast-api/store/feed"
)

type mockStore struct {
	feed  *store.Feed
	feeds []store.Feed
	err   error
}

func (m *mockStore) FindAll(ctx context.Context) ([]store.Feed, error) {
	return m.feeds, m.err
}

func (m *mockStore) FindByID(ctx context.Context, id uuid.UUID) (*store.Feed, error) {
	return m.feed, m.err
}

func (m *mockStore) FindByUserID(ctx context.Context, userID uuid.UUID) ([]store.Feed, error) {
	return m.feeds, m.err
}

func (m *mockStore) FindByIDAndUserID(ctx context.Context, id, userID uuid.UUID) (*store.Feed, error) {
	return m.feed, m.err
}

func (m *mockStore) Create(ctx context.Context, feed *store.Feed) error {
	return m.err
}

func (m *mockStore) Update(ctx context.Context, feed *store.Feed) error {
	return m.err
}

func (m *mockStore) Delete(ctx context.Context, feed *store.Feed) error {
	return m.err
}

func TestService_GetFeed(t *testing.T) {
	feed := &store.Feed{URL: "https://example.com", Title: "Example"}
	s := &mockStore{feed: feed}
	service := NewService(s)

	result, err := service.GetFeed(context.Background(), feed.ID)
	assert.NoError(t, err)
	assert.Equal(t, feed, result)
}

func TestService_GetFeed_Error(t *testing.T) {
	s := &mockStore{err: errors.New("not found")}
	service := NewService(s)

	result, err := service.GetFeed(context.Background(), uuid.Must(uuid.NewV7()))
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestService_GetFeedsByUserID(t *testing.T) {
	feeds := []store.Feed{{URL: "https://example.com", Title: "Example"}}
	s := &mockStore{feeds: feeds}
	service := NewService(s)

	result, err := service.GetFeedsByUserID(context.Background(), uuid.Must(uuid.NewV7()))
	assert.NoError(t, err)
	assert.Equal(t, feeds, result)
}

func TestService_GetFeedsByUserID_Error(t *testing.T) {
	s := &mockStore{err: errors.New("database error")}
	service := NewService(s)

	result, err := service.GetFeedsByUserID(context.Background(), uuid.Must(uuid.NewV7()))
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestService_CreateFeed(t *testing.T) {
	s := &mockStore{}
	service := NewService(s)

	feed := &store.Feed{URL: "https://example.com", Title: "Example"}
	err := service.CreateFeed(context.Background(), feed)
	assert.NoError(t, err)
}

func TestService_CreateFeed_Error(t *testing.T) {
	s := &mockStore{err: errors.New("create error")}
	service := NewService(s)

	feed := &store.Feed{URL: "https://example.com", Title: "Example"}
	err := service.CreateFeed(context.Background(), feed)
	assert.Error(t, err)
}

func TestService_DeleteFeed(t *testing.T) {
	feed := &store.Feed{URL: "https://example.com", Title: "Example"}
	s := &mockStore{feed: feed}
	service := NewService(s)

	err := service.DeleteFeed(context.Background(), uuid.Must(uuid.NewV7()), feed.ID)
	assert.NoError(t, err)
}

func TestService_DeleteFeed_NotFound(t *testing.T) {
	s := &mockStore{err: errors.New("not found")}
	service := NewService(s)

	err := service.DeleteFeed(context.Background(), uuid.Must(uuid.NewV7()), uuid.Must(uuid.NewV7()))
	assert.Error(t, err)
}

func TestService_SyncFeed(t *testing.T) {
	feed := &store.Feed{URL: "https://example.com", Title: "Example"}
	s := &mockStore{feed: feed}
	service := NewService(s)

	err := service.SyncFeed(context.Background(), uuid.Must(uuid.NewV7()), feed.ID)
	assert.NoError(t, err)
	assert.NotNil(t, feed.SyncedAt)
	assert.WithinDuration(t, time.Now(), *feed.SyncedAt, time.Second)
}

func TestService_SyncFeed_NotFound(t *testing.T) {
	s := &mockStore{err: errors.New("not found")}
	service := NewService(s)

	err := service.SyncFeed(context.Background(), uuid.Must(uuid.NewV7()), uuid.Must(uuid.NewV7()))
	assert.Error(t, err)
}
