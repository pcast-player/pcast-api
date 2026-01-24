# PCast API

[![tests](https://github.com/pcast-player/pcast-api/actions/workflows/tests.yml/badge.svg)](https://github.com/pcast-player/pcast-api/actions/workflows/tests.yml)

## What?

PCast is just a hobby project to build a little podcast player. This is the API part of it.

## Why?

Why not? But seriously, almost all big podcast player apps have sync problems with private feeds over the time and this just sucks.

So I decided to build my own podcast player. The API is written in Go Lang. The player apps will be built with Kotlin Multiplatform.

## How?

The API uses the Echo Framework, sqlc for type-safe SQL, and PostgreSQL as database.

### Installation

1. Check out the repository and install the dependencies:

```bash
mise run install
```

This will:
- Download Go dependencies
- Install swag (Swagger documentation generator)
- Install goose (database migration tool)
- Install sqlc (type-safe SQL code generator)

2. Start PostgreSQL using docker-compose:

```bash
cd docker
docker-compose up -d db
```

3. Create the database:

```bash
docker exec pcast-api-db-1 psql -U pcast -c "CREATE DATABASE pcast;"
docker exec pcast-api-db-1 psql -U pcast -c "CREATE DATABASE pcast_test;"
```

4. Run database migrations:

```bash
mise run migrate:up
```

### Running

Run the API server:

```bash
mise run run
```

Or directly:

```bash
go run -race .
```

The API will be available at http://localhost:8080.

There is also an API documentation available at http://localhost:8080/swagger/.

### Testing

Run all tests:

```bash
mise run test
```

Or with coverage:

```bash
go test -v -race -cover ./...
```

Run tests for a specific package:

```bash
go test -v -race ./store/user
go test -v -race ./integration_test/feed
```

Run a single test:

```bash
go test -v -race ./service/user -run TestService_GetUser
```

### Database Migrations

Create a new migration:

```bash
mise run migrate:create name=add_new_feature
```

Run migrations:

```bash
mise run migrate:up       # Apply all pending migrations
mise run migrate:down     # Rollback last migration
mise run migrate:status   # Check migration status
```

### Code Generation

After modifying SQL queries in `db/queries/`, regenerate the Go code:

```bash
mise run sqlc
```
