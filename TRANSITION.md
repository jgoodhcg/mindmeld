# Transition Guide: Next.js to Go Stack

This document records the transition from the original Next.js/TypeScript stack to a Go-based stack.

## Rationale

The original stack (Next.js + TypeScript + custom Node server) was replaced to achieve:

1. **Simpler deployment** - Single static binary vs Node.js runtime
2. **Type-safe templates** - templ catches HTML errors at compile time
3. **No JavaScript build step** - HTMX replaces React, eliminating bundler complexity
4. **Lower cognitive load** - Go's simplicity and explicit error handling
5. **Better WebSocket story** - nhooyr.io/websocket is more idiomatic than ws + Next.js hybrid

## What Was Removed

### Files & Directories

| Path | Description |
|------|-------------|
| `node_modules/` | NPM dependencies (~306 packages) |
| `.next/` | Next.js build output |
| `src/` | React/Next.js source code |
| `public/` | Next.js static assets (SVG icons) |
| `package.json` | NPM dependencies and scripts |
| `package-lock.json` | NPM lockfile |
| `tsconfig.json` | TypeScript configuration |
| `next.config.ts` | Next.js configuration |
| `eslint.config.mjs` | ESLint configuration |
| `postcss.config.mjs` | PostCSS configuration |
| `scripts/run-migrations.js` | Node.js migration runner |
| `db/migrations/` | Old migration format (recreated for goose) |

### Dependencies Removed

- next (16.0.8)
- react (19.2.1)
- react-dom (19.2.1)
- pg (PostgreSQL client for Node)
- typescript
- eslint + eslint-config-next
- @tailwindcss/postcss
- babel-plugin-react-compiler

## What Was Kept

| Path | Description |
|------|-------------|
| `.git/` | Git history preserved |
| `.gitignore` | Updated for Go project |
| `docker-compose.yml` | PostgreSQL setup (still valid) |
| `.env.example` | Environment template |
| `.env.local` | Local environment config |
| `DESIGN.md` | Game design document |
| `DESIGN_GAMEPLAY.md` | Gameplay overview |
| `README.md` | Updated for new stack |
| `agents.md` | Agent guidelines |

## New Project Structure

```
mindmeld/
├── cmd/server/main.go        # Entrypoint
├── internal/
│   ├── handlers/             # HTTP route handlers
│   ├── middleware/           # Auth, logging
│   ├── db/                   # sqlc generated code
│   └── ws/                   # WebSocket hub
├── templates/                # templ components (.templ)
├── static/css/input.css      # Tailwind source
├── migrations/               # goose SQL files
├── queries/                  # sqlc SQL queries
├── go.mod, go.sum            # Go dependencies
├── Makefile                  # Build orchestration
├── Dockerfile                # Container build
├── .air.toml                 # Live reload config
├── sqlc.yaml                 # sqlc config
└── tailwind.config.js        # Tailwind config
```

## Mapping: Old → New

| Old (Next.js) | New (Go) |
|---------------|----------|
| `src/app/page.tsx` | `templates/home.templ` |
| `src/app/layout.tsx` | `templates/layout.templ` |
| `src/app/globals.css` | `static/css/input.css` |
| React components | templ components |
| API routes (`/api/*`) | `internal/handlers/` |
| `scripts/run-migrations.js` | `goose` CLI |
| `npm run dev` | `make dev` (Air) |
| `npm run build` | `make build` |
| `npm run db:migrate` | `make migrate` |

## Database Schema

The schema from `db/migrations/0001_init.sql` was preserved and converted to goose format in `migrations/`. Tables remain the same:

- `game_types` - Registry of game types
- `lobbies` - Game lobbies with phases
- `teams` - Teams within lobbies
- `players` - Players with optional user_id
- `lobby_events` - Event log
- `trivia_questions` - Multiple choice questions
- `trivia_rounds` - Individual rounds
- `trivia_answers` - Player responses

## Environment Variables

| Old | New | Notes |
|-----|-----|-------|
| `DATABASE_URL` | `DATABASE_URL` | Unchanged |
| `LOBBY_TOKEN_SECRET` | `SESSION_SECRET` | Renamed for clarity |
| - | `PORT` | Server port (default 8080) |
| - | `ENV` | development/production |

## Prerequisites for New Stack

Install these tools:

```bash
# Go (1.22+)
brew install go

# templ
go install github.com/a-h/templ/cmd/templ@latest

# Air (live reload)
go install github.com/air-verse/air@latest

# goose (migrations)
go install github.com/pressly/goose/v3/cmd/goose@latest

# sqlc
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Tailwind CSS standalone CLI
curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-macos-arm64
chmod +x tailwindcss-macos-arm64
mv tailwindcss-macos-arm64 /usr/local/bin/tailwindcss
```

## Getting Started After Transition

```bash
# Start PostgreSQL
docker compose up -d

# Run migrations
make migrate

# Generate code (templ + sqlc)
make generate

# Start development server with live reload
make dev
```

## References

- [templ documentation](https://templ.guide)
- [HTMX documentation](https://htmx.org/docs/)
- [goose migrations](https://github.com/pressly/goose)
- [sqlc documentation](https://sqlc.dev)
- [Air live reload](https://github.com/air-verse/air)
- [chi router](https://github.com/go-chi/chi)
- [scs sessions](https://github.com/alexedwards/scs)
