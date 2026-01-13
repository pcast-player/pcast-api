# AGENTS.md - Coding Agent Guide for PCast API

This guide provides essential information for AI coding agents working in this repository.

## Project Overview

PCast API is a podcast player backend written in Go 1.22+ using the Echo framework, GORM ORM, and SQLite/PostgreSQL databases. The codebase follows a clean 3-layer architecture: Controller → Service → Store.

## Build, Test, and Run Commands

### Installation
```bash
make install          # Download dependencies and install swag tool
go mod download       # Download dependencies only
```

### Building
```bash
make build           # Build with race detector to target/
go build -o target/ -race .
```

### Running
```bash
make run             # Run with race detector
go run -race .       # Alternative: run directly
```

The API starts at http://localhost:8080 with Swagger docs at http://localhost:8080/swagger/

### Testing

```bash
# Run all tests
make test            # Run all tests with race detector
go test -v -race ./...

# Run tests in a specific package
go test -v -race ./service/user
go test -v -race ./store/feed
go test -v -race ./integration_test/feed

# Run a single test
go test -v -race ./service/user -run TestService_GetUser
go test -v -race ./store/feed -run TestStore_Create
```

### Documentation
```bash
make create-docs     # Generate Swagger docs with swag init
```

## Project Structure

```
pcast-api/
├── controller/          # HTTP handlers, request/response DTOs, validation
│   ├── service_interface/   # Service interfaces for DI
│   └── {domain}/           # handler.go, presenter.go, *_request.go
├── service/             # Business logic layer
│   ├── model_interface/     # Store interfaces for DI
│   └── {domain}/           # service.go, service_test.go
├── store/               # Data access layer
│   └── {domain}/           # model.go, store.go, store_test.go
├── router/              # Echo router setup
├── db/                  # Database connection
├── config/              # TOML configuration
├── helper/              # Test utilities
└── integration_test/    # Full-stack API tests
```

## Code Style Guidelines

### Import Organization

Always use this three-group pattern with blank lines between groups:

```go
import (
    // 1. Standard library
    "fmt"
    "time"
    
    // 2. Third-party packages
    "github.com/google/uuid"
    "github.com/labstack/echo/v4"
    "gorm.io/gorm"
    
    // 3. Internal packages (pcast-api/*)
    "pcast-api/config"
    "pcast-api/store/feed"
)
```

### Import Aliasing

Use aliases to avoid conflicts and improve clarity:

```go
import (
    serviceInterface "pcast-api/controller/service_interface"
    modelInterface "pcast-api/service/model_interface"
    feedStore "pcast-api/store/feed"
    userStore "pcast-api/store/user"
)
```

### Naming Conventions

- **Files**: lowercase with underscores: `handler.go`, `service_test.go`, `create_request.go`
- **Constructors**: Use `New()` for main type, `NewXxx()` for specific types: `NewHandler()`, `NewPresenter()`
- **Request DTOs**: Suffix with `Request`: `CreateRequest`, `UpdatePasswordRequest`
- **Response DTOs**: Name as `Presenter` with constructor: `type Presenter struct` + `func NewPresenter()`
- **Test functions**: `TestStructName_MethodName` or `TestMethodName`
- **Variables**: camelCase for locals, PascalCase for exported

### Type Definitions

**Models (Store Layer):**
```go
type Feed struct {
    ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
    CreatedAt time.Time
    UpdatedAt time.Time
    UserID    uuid.UUID `gorm:"type:uuid"`
    Title     string
    URL       string
    SyncedAt  *time.Time  // Use pointer for nullable fields
}

func (f *Feed) BeforeCreate(_ *gorm.DB) (err error) {
    f.ID, err = uuid.NewV7()  // Always use UUID v7
    return err
}
```

**Request DTOs (Controller Layer):**
```go
type CreateRequest struct {
    Title string `json:"title" validate:"required"`
    URL   string `json:"url" validate:"required,url"`
}
```

**Response DTOs (Controller Layer):**
```go
type Presenter struct {
    ID       uuid.UUID  `json:"id"`
    Title    string     `json:"title"`
    SyncedAt *time.Time `json:"syncedAt"`
}

func NewPresenter(feed *feed.Feed) *Presenter {
    return &Presenter{
        ID:       feed.ID,
        Title:    feed.Title,
        SyncedAt: feed.SyncedAt,
    }
}
```

### Struct and Constructor Patterns

Every major struct has a constructor:

```go
type Handler struct {
    service serviceInterface.Feed
}

func NewHandler(service serviceInterface.Feed) *Handler {
    return &Handler{service: service}
}
```

Handlers have a `Register()` method for mounting routes:

```go
func (h *Handler) Register(g *echo.Group) {
    g.GET("/feeds", h.GetFeeds)
    g.POST("/feeds", h.CreateFeed)
    g.DELETE("/feeds/:id", h.DeleteFeed)
}
```

### Error Handling

Keep error handling simple and direct:

```go
func (h *Handler) GetFeeds(c echo.Context) error {
    userID, err := getUser(c)
    if err != nil {
        return c.NoContent(http.StatusUnauthorized)
    }
    
    feeds, err := h.service.GetFeedsByUserID(userID)
    if err != nil {
        return c.NoContent(http.StatusInternalServerError)
    }
    
    // Transform and return
    res := lo.Map(feeds, func(item model.Feed, index int) *Presenter {
        return NewPresenter(&item)
    })
    
    return c.JSON(http.StatusOK, res)
}
```

- Return errors immediately with appropriate HTTP status
- Use `c.NoContent(status)` for errors without body
- Use `c.JSON(status, data)` for success responses
- Propagate GORM errors from store through service to controller

### Swagger Documentation

Add godoc comments with Swagger annotations to all handlers:

```go
// GetFeeds godoc
// @Summary Get all feeds
// @Description Retrieve all feeds for the authenticated user
// @Tags feeds
// @Produce json
// @Param Authorization header string true "User ID"
// @Success 200 {array} Presenter
// @Router /feeds [get]
func (h *Handler) GetFeeds(c echo.Context) error {
    // implementation
}
```

## Testing Patterns

### Unit Tests (Service Layer)

Create manual mocks implementing interfaces:

```go
type mockStore struct {
    user *store.User
    err  error
}

func (m *mockStore) FindByID(id uuid.UUID) (*store.User, error) {
    return m.user, m.err
}

func TestService_GetUser(t *testing.T) {
    user := &store.User{Email: "foo@bar.com"}
    s := &mockStore{user: user}
    service := NewService(s)
    
    result, err := service.GetUser(user.ID)
    assert.NoError(t, err)
    assert.Equal(t, user, result)
}
```

### Store Tests

Use `TestMain` for setup/teardown with SQLite:

```go
func TestMain(m *testing.M) {
    setup()
    code := m.Run()
    tearDown()
    os.Exit(code)
}
```

### Integration Tests

Use `apitest` library for full HTTP testing:

```go
func TestCreateFeed(t *testing.T) {
    userID := createUser(t)
    
    apitest.New().
        Handler(newApp()).
        Post("/api/feeds").
        Header("Authorization", userID.String()).
        JSON(`{"url": "https://example.com","title":"Example"}`).
        Expect(t).
        Status(http.StatusCreated).
        Assert(jsonpath.Equal("$.title", "Example")).
        End()
    
    truncateTables()
}
```

## Additional Notes

- **UUID v7**: Always use `uuid.NewV7()` for ID generation (not v4)
- **Password Hashing**: Use `argon2id.CreateHash()` and `argon2id.ComparePasswordAndHash()`
- **Functional Utils**: Use `samber/lo` for Map, Filter, etc.
- **Database**: GORM with AutoMigrate in store constructors
- **Auth**: Currently simple UUID in Authorization header (no JWT middleware yet)
- **Config**: TOML-based configuration in `config.toml`

## Common Pitfalls

- Don't forget blank lines between import groups
- Always use pointer receivers for struct methods
- Use pointer types for nullable database fields
- Call `truncateTables()` after each integration test
- Enable race detector in tests: `-race` flag
- Don't forget to regenerate Swagger docs after API changes
