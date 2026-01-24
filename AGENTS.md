# AGENTS.md - Coding Agent Guide for PCast API

This guide provides essential information for AI coding agents working in this repository.

## Project Overview

PCast API is a podcast player backend written in Go 1.22+ using the Echo framework, sqlc for type-safe SQL, and PostgreSQL database. The codebase follows a clean 3-layer architecture: Controller → Service → Store.

## Quick Reference

```bash
# Common commands
mise run install         # Install dependencies, swag, goose, sqlc
mise run build           # Build with race detector
mise run run             # Run the API server
mise run test            # Run all tests with race detector
mise run docs            # Generate Swagger documentation

# Database migrations
mise run migrate:up      # Run database migrations
mise run migrate:status  # Check migration status
mise run sqlc            # Generate sqlc code

# Run a single test
go test -v -race ./service/user -run TestService_GetUser

# Run tests with coverage
go test -v -race -cover ./...
```

## Build, Test, and Run Commands

### Installation
```bash
mise run install      # Download dependencies and install tools (swag, goose, sqlc)
go mod download       # Download dependencies only
```

**Required Tools:**
- `swag` - Swagger documentation generator
- `goose` - Database migration tool
- `sqlc` - Type-safe SQL code generator

### Building
```bash
mise run build       # Build with race detector to target/
go build -o target/ -race .
```

### Running
```bash
mise run run         # Run with race detector
go run -race .       # Alternative: run directly
```

The API starts at http://localhost:8080 with Swagger docs at http://localhost:8080/swagger/

### Testing

```bash
# Run all tests
mise run test        # Run all tests with race detector
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
mise run docs        # Generate Swagger docs with swag init
```

### Database Migrations
```bash
# Run migrations
mise run migrate:up              # Apply all pending migrations
mise run migrate:down            # Rollback last migration
mise run migrate:status          # Check migration status

# Create new migration
mise run migrate:create name=add_new_feature

# Test database migrations
mise run migrate:test:up         # Apply migrations to test database
mise run migrate:test:down       # Rollback test database migrations
```

### Generate sqlc Code
```bash
mise run sqlc         # Generate type-safe Go code from SQL queries
sqlc generate         # Alternative: run directly
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
├── db/                  # Database connection & sqlc generated code
│   ├── migrations/         # Goose SQL migrations
│   ├── queries/            # sqlc SQL query definitions
│   ├── sqlcgen/            # Generated sqlc code (DO NOT EDIT)
│   └── db.go               # Database connection setup
├── config/              # TOML configuration
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
    ID        uuid.UUID
    CreatedAt time.Time
    UpdatedAt time.Time
    UserID    uuid.UUID
    Title     string
    URL       string
    SyncedAt  *time.Time  // Use pointer for nullable fields
}

// BeforeCreate sets default values before creating a feed
// Call this explicitly in Store.Create()
func (f *Feed) BeforeCreate() error {
    if f.ID == uuid.Nil {
        id, err := uuid.NewV7()  // Always use UUID v7
        if err != nil {
            return err
        }
        f.ID = id
    }
    
    if f.CreatedAt.IsZero() {
        f.CreatedAt = time.Now()
    }
    
    if f.UpdatedAt.IsZero() {
        f.UpdatedAt = time.Now()
    }
    
    return nil
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

Use `TestMain` for setup/teardown with PostgreSQL:

```go
var d *sql.DB
var fs *Store

const testDSN = "host=localhost port=5432 user=pcast password=pcast dbname=pcast_test sslmode=disable"

func TestMain(m *testing.M) {
    setup()
    code := m.Run()
    tearDown()
    os.Exit(code)
}

func setup() {
    d = db.NewTestDB(testDSN)
    runMigrations()
    fs = New(d)
}

func tearDown() {
    truncateTable()
    d.Close()
}

func runMigrations() {
    _, err := d.Exec(`
        CREATE TABLE IF NOT EXISTS feeds (
            id UUID PRIMARY KEY,
            created_at TIMESTAMP NOT NULL DEFAULT NOW(),
            updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
            user_id UUID NOT NULL,
            title VARCHAR(500) NOT NULL,
            url VARCHAR(1000) NOT NULL,
            synced_at TIMESTAMP
        );
        CREATE INDEX IF NOT EXISTS idx_feeds_user_id ON feeds(user_id);
    `)
    if err != nil {
        panic(fmt.Sprintf("failed to run migrations: %v", err))
    }
}

func truncateTable() {
    d.Exec("TRUNCATE TABLE feeds")
}

func TestCreateFeed(t *testing.T) {
    userID, _ := uuid.NewV7()
    feed := &Feed{URL: "https://example.com", Title: "Example", UserID: userID}
    err := fs.Create(feed)
    assert.NoError(t, err)
    
    truncateTable()  // Clean up after each test
}
```

**Store test patterns:**
- All tests use PostgreSQL (`pcast_test` database)
- Call `truncateTable()` after each test to clean data
- Use direct SQL `TRUNCATE TABLE` for cleanup
- `runMigrations()` creates tables in setup

### Integration Tests

Use `apitest` library for full HTTP testing:

```go
var sqlDB *sql.DB

const testDSN = "host=localhost port=5432 user=pcast password=pcast dbname=pcast_test sslmode=disable"

func TestMain(m *testing.M) {
    sqlDB = db.NewTestDB(testDSN)
    code := m.Run()
    sqlDB.Exec("TRUNCATE TABLE users CASCADE")
    sqlDB.Exec("TRUNCATE TABLE feeds CASCADE")
    sqlDB.Close()
    os.Exit(code)
}

func truncateTables() {
    sqlDB.Exec("TRUNCATE TABLE users CASCADE")
    sqlDB.Exec("TRUNCATE TABLE feeds")
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
- All tests use PostgreSQL (`pcast_test` database)
- Use `apitest` for HTTP testing and `apitest-jsonpath` for assertions
- Call `truncateTables()` after each test
- Use CASCADE for truncate to handle foreign keys
- Helper functions like `createUser()` for test data setup

## Additional Notes

- **UUID v7**: Always use `uuid.NewV7()` for ID generation (not v4)
- **Password Hashing**: Use `argon2id.CreateHash()` and `argon2id.ComparePasswordAndHash()`
- **Functional Utils**: Use `samber/lo` for Map, Filter, etc.
- **Database**: sqlc for type-safe SQL queries with PostgreSQL
- **Migrations**: Goose for SQL migrations in `db/migrations/`
- **Auth**: Currently simple UUID in Authorization header (no JWT middleware yet)
- **Config**: TOML-based configuration in `config.toml`
- **Validation**: Custom validator using `go-playground/validator` configured in router
- **Service Layer**: Often thin pass-through layer coordinating between controllers and stores
- **sqlc Generated Code**: Located in `db/sqlcgen/` (DO NOT manually edit)
- **Ignored Files**: `.gitignore` excludes `*.db`, `target/`, and `.idea/`
- **Conventional Commits**: Use the [Conventional Commits](https://www.conventionalcommits.org/) standard for commit messages

## File Organization

### Test Database
- All tests use PostgreSQL (`pcast_test` database)
- Tests connect via: `host=localhost port=5432 user=pcast password=pcast dbname=pcast_test sslmode=disable`
- Use docker-compose to start Postgres: `cd docker && docker-compose up -d db`
- Create test database: `docker exec pcast-api-db-1 psql -U pcast -c "CREATE DATABASE pcast_test;"`

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
  model.go            # Domain model with BeforeCreate() method
  store.go            # Data access methods using sqlc
  store_test.go       # Store tests with PostgreSQL

db/
  migrations/         # Goose SQL migration files
  queries/            # sqlc SQL query definitions
  sqlcgen/            # Generated sqlc code (DO NOT EDIT)
```

## Working with sqlc

### Adding New Queries
1. Write SQL in `db/queries/{domain}.sql` with sqlc annotations
2. Run `mise run sqlc` to generate Go code
3. Use generated code in store layer

**Example query:**
```sql
-- name: FindUserByEmail :one
SELECT * FROM users WHERE email = $1;
```

**Generated usage:**
```go
user, err := s.queries.FindUserByEmail(context.Background(), email)
```

### Creating Migrations
```bash
# Create new migration
mise run migrate:create name=add_user_avatar

# Edit the generated SQL file in db/migrations/
# Run migration
mise run migrate:up
```

### Store Layer Pattern with sqlc
```go
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

func (s *Store) Create(feed *Feed) error {
    if err := feed.BeforeCreate(); err != nil {
        return err
    }
    
    _, err := s.queries.CreateFeed(context.Background(), sqlcgen.CreateFeedParams{
        ID:        feed.ID,
        CreatedAt: feed.CreatedAt,
        // ... other fields
    })
    
    return err
}
```

## Common Pitfalls

- Don't forget blank lines between import groups
- Always use pointer receivers for struct methods
- Use pointer types for nullable database fields
- Call `truncateTable()` after each test to clean data
- Enable race detector in tests: `-race` flag
- Don't forget to regenerate Swagger docs after API changes: `mise run docs`
- Import domain stores with appropriate aliases to avoid conflicts
- **sqlc**: Don't manually edit files in `db/sqlcgen/` - always regenerate
- **Migrations**: Always create both Up and Down migrations
- **NULL handling**: Use sql.NullTime, sql.NullInt32, etc. for nullable fields
- Call `BeforeCreate()` explicitly in Store.Create() before inserting
