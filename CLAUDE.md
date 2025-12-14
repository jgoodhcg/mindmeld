# Agent Guidelines

## Development

- Do not start the dev server automatically; let the user run it to avoid port conflicts.
- Dev tools (air, templ, sqlc, goose) are pinned in go.mod - use `go run` instead of global installs.
- Generated files (templates/*_templ.go, internal/db/, static/css/output.css) are not committed - run `make generate` and `make css` locally.

## Stack

- Go server with chi router
- templ for type-safe HTML templates
- Tailwind CSS (standalone CLI)
- PostgreSQL with sqlc for type-safe queries
- goose for migrations

## Key Commands

```bash
make dev        # Start with hot reload
make generate   # Generate templ + sqlc code
make css        # Build Tailwind CSS
make migrate    # Run database migrations
```
