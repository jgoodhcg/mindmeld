.PHONY: dev build run templ css css-watch clean db-up db-down

# Load environment variables
include .env.local
export

# Development with live reload (pinned version via tools.go)
dev:
	go run github.com/air-verse/air

# Build production binary
build: templ css
	go build -o bin/server ./cmd/server

# Run the server directly
run: templ css
	go run ./cmd/server

# Generate templ templates (pinned version via tools.go)
templ:
	go run github.com/a-h/templ/cmd/templ generate

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

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f static/css/output.css
	rm -f templates/*_templ.go
