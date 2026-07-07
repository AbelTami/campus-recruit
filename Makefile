.PHONY: run build test clean migrate-up migrate-down seed

# Load .env file if exists
-include .env
export

run:
	go run cmd/api/main.go

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/server cmd/api/main.go

test:
	go test -v -race -coverprofile=coverage.out ./...

clean:
	rm -rf bin/ coverage.out

migrate-up:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" up

migrate-down:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" down

seed:
	go run cmd/seed/main.go

lint:
	golangci-lint run ./...

mod:
	go mod tidy

dev: mod run
