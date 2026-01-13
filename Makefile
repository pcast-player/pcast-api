all: install build fixtures

install:
	go mod download
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

build:
	go build -o target/ -race .

test:
	go test -v -race ./...

run:
	go run -race .

create-docs:
	swag init --parseDependency --parseInternal

# Database migration commands
DB_DSN ?= "host=localhost port=5432 user=pcast password=pcast dbname=pcast sslmode=disable"
TEST_DB_DSN ?= "host=localhost port=5432 user=pcast password=pcast dbname=pcast_test sslmode=disable"

migrate-up:
	cd db/migrations && goose postgres $(DB_DSN) up

migrate-down:
	cd db/migrations && goose postgres $(DB_DSN) down

migrate-status:
	cd db/migrations && goose postgres $(DB_DSN) status

migrate-create:
	@test -n "$(name)" || (echo "Error: name is required. Use: make migrate-create name=your_migration_name" && exit 1)
	cd db/migrations && goose create $(name) sql

test-migrate-up:
	cd db/migrations && goose postgres $(TEST_DB_DSN) up

test-migrate-down:
	cd db/migrations && goose postgres $(TEST_DB_DSN) down

# Generate sqlc code
sqlc-generate:
	sqlc generate
