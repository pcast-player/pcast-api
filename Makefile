all: install build test

install:
	go mod download
	go install github.com/swaggo/swag/cmd/swag@latest

build:
	go build -o target/ -race .

test:
	go test -v -race ./...

run:
	go run -race .

create-docs:
	swag init --parseDependency --parseInternal
