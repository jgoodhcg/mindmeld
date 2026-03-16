.PHONY: dev build run generate templ sqlc css clean db-up db-down migrate migrate-down migrate-status fmt lint cluster-content-validate cluster-content-import e2e-install e2e-screenshot e2e-flow e2e-multiplayer e2e-test

# Source .env.local with shell semantics so `make` matches `source .env.local && ...`.
WITH_DOTENV = if [ -f .env.local ]; then set -a; . ./.env.local; set +a; fi;

# Development with live reload
dev:
	@$(WITH_DOTENV) go run github.com/air-verse/air

# Build production binary
build: generate css
	@$(WITH_DOTENV) go build -o bin/server ./cmd/server

# Run the server directly
run: generate css
	@$(WITH_DOTENV) go run ./cmd/server

# Generate all code (templ + sqlc)
generate: templ sqlc

# Format all code (templ + go)
fmt:
	go run github.com/a-h/templ/cmd/templ fmt .
	go fmt ./...

# Lint code (go vet)
lint:
	go vet ./...

# Validate cluster content library JSON without DB writes
cluster-content-validate:
	@$(WITH_DOTENV) go run ./cmd/cluster-content validate -file content/cluster/library.v1.json

# Import cluster content library into database
cluster-content-import:
	@$(WITH_DOTENV) go run ./cmd/cluster-content import -file content/cluster/library.v1.json

# Generate templ templates
templ:
	go run github.com/a-h/templ/cmd/templ generate

# Generate sqlc code
sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc generate

# Build Tailwind CSS
css:
	tailwindcss -i styles/input.css -o static/css/output.css --minify

# Start database
db-up:
	docker compose up -d

# Stop database
db-down:
	docker compose down

# Run migrations
migrate:
	@$(WITH_DOTENV) go run github.com/pressly/goose/v3/cmd/goose -dir migrations postgres "$${DATABASE_URL}" up

# Rollback last migration
migrate-down:
	@$(WITH_DOTENV) go run github.com/pressly/goose/v3/cmd/goose -dir migrations postgres "$${DATABASE_URL}" down

# Migration status
migrate-status:
	@$(WITH_DOTENV) go run github.com/pressly/goose/v3/cmd/goose -dir migrations postgres "$${DATABASE_URL}" status

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f static/css/output.css
	rm -f templates/*_templ.go
	rm -rf internal/db/

# E2E: Install playwright and browsers
e2e-install:
	cd e2e && npm install && npm run install-browsers

# E2E: Take a screenshot (usage: make e2e-screenshot or make e2e-screenshot ARGS="/lobby/ABC123")
e2e-screenshot:
	cd e2e && npm run screenshot -- $(ARGS)

# E2E: Run a UI flow (usage: make e2e-flow or make e2e-flow ARGS="join ABC123")
e2e-flow:
	cd e2e && npm run flow -- $(ARGS)

# E2E: Run multiplayer round simulation
e2e-multiplayer:
	cd e2e && npm run multiplayer

# E2E: Run playwright tests
e2e-test:
	cd e2e && npm test
