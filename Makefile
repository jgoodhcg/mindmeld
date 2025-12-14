.PHONY: dev build run generate templ sqlc css css-watch clean db-up db-down migrate migrate-down migrate-status

# Load environment variables
include .env.local
export

# Development with live reload
dev:
	go run github.com/air-verse/air

# Build production binary
build: generate css
	go build -o bin/server ./cmd/server

# Run the server directly
run: generate css
	go run ./cmd/server

# Generate all code (templ + sqlc)
generate: templ sqlc

# Generate templ templates
templ:
	go run github.com/a-h/templ/cmd/templ generate

# Generate sqlc code
sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc generate

# Build Tailwind CSS
css:
	tailwindcss -i static/css/input.css -o static/css/output.css --minify

# Watch Tailwind CSS
css-watch:
	tailwindcss -i static/css/input.css -o static/css/output.css --watch

# Start database
db-up:
	docker compose up -d

# Stop database
db-down:
	docker compose down

# Run migrations
migrate:
	go run github.com/pressly/goose/v3/cmd/goose -dir migrations postgres "$(DATABASE_URL)" up

# Rollback last migration
migrate-down:
	go run github.com/pressly/goose/v3/cmd/goose -dir migrations postgres "$(DATABASE_URL)" down

# Migration status
migrate-status:
	go run github.com/pressly/goose/v3/cmd/goose -dir migrations postgres "$(DATABASE_URL)" status

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f static/css/output.css
	rm -f templates/*_templ.go
	rm -rf internal/db/
