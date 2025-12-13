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
- [Tailwind CSS CLI](https://github.com/tailwindlabs/tailwindcss/releases) (standalone binary)

### Install Tailwind CLI (macOS ARM)

```bash
curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-macos-arm64
chmod +x tailwindcss-macos-arm64
sudo mv tailwindcss-macos-arm64 /usr/local/bin/tailwindcss
```

## Getting Started

```bash
# Copy environment config
cp .env.example .env.local

# Start PostgreSQL
make db-up

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
| `make clean` | Remove build artifacts |

## Project Structure

```
mindmeld/
├── cmd/server/main.go    # Server entrypoint
├── templates/            # templ templates
├── static/css/           # Tailwind CSS
├── tools.go              # Pinned tool versions
├── Makefile              # Build commands
├── .air.toml             # Hot reload config
└── docker-compose.yml    # PostgreSQL
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `DATABASE_URL` | PostgreSQL connection string |
| `PORT` | Server port (default: 3000) |
