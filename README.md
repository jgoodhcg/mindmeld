# Mindmeld

Multiplayer party games that bring people together through shared thinking.

## Stack

- **Go** - Server
- **templ** - Type-safe HTML templates
- **HTMX** - Frontend interactivity (coming soon)
- **Tailwind CSS** - Styling
- **PostgreSQL** - Database
- **Air** - Hot reload for development

## Prerequisites

- Go 1.25+
- Docker (for PostgreSQL)
- [Tailwind CSS CLI](https://github.com/tailwindlabs/tailwindcss/releases/tag/v4.1.18) - download the standalone binary for your platform and add to PATH as `tailwindcss`

## Getting Started

```bash
# Copy environment config
cp .env.example .env.local

# Start PostgreSQL
make db-up

# Run migrations
make migrate

# Build CSS
make css

# Start dev server with hot reload
make dev
```

Open [http://localhost:3000](http://localhost:3000) to see the app.

## Available Commands

| Command | Description |
|---------|-------------|
| `make dev` | Start server with hot reload |
| `make run` | Run server directly |
| `make build` | Build production binary |
| `make templ` | Generate templ templates |
| `make css` | Build Tailwind CSS |
| `make css-watch` | Watch and rebuild CSS |
| `make db-up` | Start PostgreSQL container |
| `make db-down` | Stop PostgreSQL container |
| `make migrate` | Run database migrations |
| `make clean` | Remove build artifacts |

## Project Structure

```
mindmeld/
├── cmd/server/main.go    # Server entrypoint
├── templates/            # templ templates (*.templ source files)
├── static/css/           # Tailwind CSS
├── tools.go              # Pinned tool versions
├── Makefile              # Build commands
├── .air.toml             # Hot reload config
├── docker-compose.yml    # Local dev PostgreSQL
└── Dockerfile            # Production build
```

## Docker Files

| File | Purpose | When to use |
|------|---------|-------------|
| `docker-compose.yml` | Local Postgres only | `make db-up` for local dev with hot reload |
| `Dockerfile` | Production app build | DO App Platform, or test locally with `docker build` |

**Why separate?** Local dev runs Go directly with Air for hot reload. The Dockerfile is for production builds. Combining them would sacrifice hot reload for no benefit.

## Environment Variables

| Variable | Description |
|----------|-------------|
| `DATABASE_URL` | PostgreSQL connection string |
| `PORT` | Server port (default: 8080) |

## Deployment

### Digital Ocean App Platform

Deploys via Dockerfile. The Dockerfile handles all code generation (templ, sqlc, Tailwind) at build time, so generated files are not committed to git.

**Test the Docker build locally:**

```bash
# Build the image
docker build -t mindmeld .

# Run it (with local Postgres via docker-compose)
make db-up
docker run -p 8080:8080 -e DATABASE_URL="postgres://mindmeld:mindmeld@host.docker.internal:5432/mindmeld" mindmeld
```

**DO App Platform setup:**
1. Connect your GitHub repo
2. DO auto-detects the Dockerfile
3. Add a managed Postgres database
4. Set `DATABASE_URL` env var (auto-injected from managed Postgres)
5. HTTP Port: `8080`

### Migrations:

Run migrations against production DB after deploy (run goose directly to avoid `.env.local` override):
```bash
go run github.com/pressly/goose/v3/cmd/goose -dir migrations postgres "your-prod-connection-string" up
```
