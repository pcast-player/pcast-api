# PCast API

[![tests](https://github.com/pcast-player/pcast-api/actions/workflows/tests.yml/badge.svg)](https://github.com/pcast-player/pcast-api/actions/workflows/tests.yml)

## What?

PCast is just a hobby project to build a little podcast player. This is the API part of it.

## Why?

Why not? But seriously, almost all big podcast player apps have sync problems with private feeds over the time and this just sucks.

So I decided to build my own podcast player. The API is written in Go Lang. The player apps will be built with Kotlin Multiplatform.

## How?

The API uses the Echo Framework and SQLite as database (for now).

### Installation

Check out the repository and install the dependencies with go mod:

```bash
go mod download
```

### Running

Just run the main.go file:

```bash
go run main.go
```

The API will be available at http://localhost:8080.

### Testing

```bash
go test ./...
```
