# AGENTS.md - Coding Agent Guide for PCast API

This guide provides essential information for AI coding agents working in this repository.

## Project Overview

PCast API is a podcast player backend written in Go 1.22+ using the Echo framework, GORM ORM, and SQLite/PostgreSQL databases. The codebase follows a clean 3-layer architecture: Controller → Service → Store.

## Quick Reference

```bash
# Common commands
make install         # Install dependencies & swag tool
make build           # Build with race detector
make run             # Run the API server
make test            # Run all tests with race detector
make create-docs     # Generate Swagger documentation

# Run a single test
go test -v -race ./service/user -run TestService_GetUser

# Run tests with coverage
go test -v -race -cover ./...
```

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

# Run tests with coverage
go test -v -race -cover ./...

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

**In controllers**, when importing domain models for type usage, use `model` alias:
```go
import (
    model "pcast-api/store/feed"
)
```

### Naming Conventions

- **Files**: lowercase with underscores: `handler.go`, `service_test.go`, `create_request.go`
- **Constructors**: Use `New()` for main type, `NewXxx()` for specific types: `NewHandler()`, `NewPresenter()`
- **Request DTOs**: Suffix with `Request`: `CreateRequest`, `UpdatePasswordRequest`
- **Response DTOs**: Name as `Presenter` with constructor: `type Presenter struct` + `func NewPresenter()`
- **Test functions**: `TestStructName_MethodName` or `TestMethodName`
- **Variables**: camelCase for locals, PascalCase for exported
- **Handler methods**: Can be unexported (lowercase) if only called via `Register()`: `registerUser()`, `updatePassword()`
- **Request variables**: Use short names: `r := new(CreateRequest)` or `userRequest := new(RegisterRequest)`

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
// CreateRequest represents a feed request
// @model CreateRequest
type CreateRequest struct {
    Title string `json:"title" validate:"required"`
    URL   string `json:"url" validate:"required,url"`
}
```

Common validation tags:
- `validate:"required"` - field is required
- `validate:"email"` - valid email format
- `validate:"url"` - valid URL format

**Response DTOs (Controller Layer):**
```go
// Presenter represents a feed presenter
// @model Presenter
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

Note: Presenters filter out sensitive data (e.g., passwords are never included in User presenters).

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

Note: Only the `Register()` method needs to be exported; handler methods can be unexported.

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

**Error handling patterns:**
- Return errors immediately with appropriate HTTP status
- Use `c.NoContent(status)` for errors without body
- Use `c.JSON(status, data)` for success responses
- Propagate GORM errors from store through service to controller
- For `c.Bind()` and `c.Validate()` errors, return the error directly (Echo handles it)

**Request binding and validation:**
```go
r := new(CreateRequest)
if err := c.Bind(r); err != nil {
    return err  // Echo handles validation errors
}
if err := c.Validate(r); err != nil {
    return err  // Returns 400 with validation details
}
```

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
var d *gorm.DB
var fs *Store

func TestMain(m *testing.M) {
    setup()
    code := m.Run()
    tearDown()
    os.Exit(code)
}

func setup() {
    d = db.NewTestDB("./../../fixtures/test/store_feed.db")
    fs = New(d)
}

func tearDown() {
    helper.RemoveTable(d, &Feed{})
}

func truncateTable() {
    helper.TruncateTables(d, "feeds")
}

func TestCreateFeed(t *testing.T) {
    feed := &Feed{URL: "https://example.com"}
    err := fs.Create(feed)
    assert.NoError(t, err)
    
    truncateTable()  // Clean up after each test
}
```

**Store test patterns:**
- Test database files: `./../../fixtures/test/{package_name}.db`
- Call `truncateTable()` after each test to clean data
- Use `helper.TruncateTables(db, "table_name")` for cleanup
- Use `helper.RemoveTable(db, &Model{})` in teardown

### Integration Tests

Use `apitest` library for full HTTP testing:

```go
var d *gorm.DB

func TestMain(m *testing.M) {
    d = db.NewTestDB("./../../fixtures/test/integration_feed.db")
    code := m.Run()
    helper.RemoveTable(d, &feedStore.Feed{})
    helper.RemoveTable(d, &userStore.User{})
    os.Exit(code)
}

func truncateTables() {
    helper.TruncateTables(d, "feeds")
    helper.TruncateTables(d, "users")
}

func createUser(t *testing.T) uuid.UUID {
    result := apitest.New().
        Handler(newApp()).
        Post("/api/user/register").
        JSON(`{"email": "foo@bar.com", "password": "test"}`).
        Expect(t).
        Status(http.StatusCreated).
        End()
    
    u := unmarshal[user.Presenter](t, &result)
    return u.ID
}

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
    
    truncateTables()  // Always clean up after each test
}
```

**Integration test patterns:**
- Test database files: `./../../fixtures/test/integration_{domain}.db`
- Use `apitest` for HTTP testing and `apitest-jsonpath` for assertions
- Call `truncateTables()` after each test
- Helper functions like `createUser()` for test data setup

## Additional Notes

- **UUID v7**: Always use `uuid.NewV7()` for ID generation (not v4)
- **Password Hashing**: Use `argon2id.CreateHash()` and `argon2id.ComparePasswordAndHash()`
- **Functional Utils**: Use `samber/lo` for Map, Filter, etc.
- **Database**: GORM with AutoMigrate in store constructors
- **Auth**: Currently simple UUID in Authorization header (no JWT middleware yet)
- **Config**: TOML-based configuration in `config.toml`
- **Validation**: Custom validator using `go-playground/validator` configured in router
- **Service Layer**: Often thin pass-through layer coordinating between controllers and stores
- **Test Helpers**: `helper` package provides `TruncateTables()` and `RemoveTable()` utilities
- **Ignored Files**: `.gitignore` excludes `*.db`, `target/`, and `.idea/`

## File Organization

### Test Database Files
- Store tests: `./../../fixtures/test/store_{domain}.db`
- Integration tests: `./../../fixtures/test/integration_{domain}.db`
- All `.db` files are gitignored

### Domain Structure
Each domain (user, feed, episode) follows the same pattern:
```
controller/{domain}/
  handler.go          # HTTP handlers with Register() method
  presenter.go        # Response DTO with NewPresenter()
  *_request.go        # Request DTOs (create_request.go, etc.)

service/{domain}/
  service.go          # Business logic
  service_test.go     # Unit tests with mocks

store/{domain}/
  model.go            # GORM model with BeforeCreate hook
  store.go            # Data access methods
  store_test.go       # Store tests with SQLite
```

## Common Pitfalls

- Don't forget blank lines between import groups
- Always use pointer receivers for struct methods
- Use pointer types for nullable database fields
- Call `truncateTable()` or `truncateTables()` after each test
- Enable race detector in tests: `-race` flag
- Don't forget to regenerate Swagger docs after API changes: `make create-docs`
- Test database paths use `./../../fixtures/test/` relative to test file location
- Import domain stores with appropriate aliases to avoid conflicts
