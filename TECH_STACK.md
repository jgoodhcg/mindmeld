# Technical Architecture: Mindmeld (Go Stack)

## 1. Core Stack

- **Language:** Go (Golang)
  - *Why:* Simplicity, stability ("boring" in a good way), exceptional standard library, easy single-binary deployment. Low cognitive load.

- **Templating:** templ (https://templ.guide)
  - *Why:* Type-safe HTML. Catches errors at compile time, preventing runtime bugs in UI code.

- **Frontend Logic:** HTMX + WebSockets
  - *Why:* HATEOAS architecture. No frontend build step, no JavaScript framework complexity.

- **Styling:** Tailwind CSS (CLI)
  - *Why:* Rapid UI development. The standalone CLI means no Node.js/NPM dependency is required for builds.

## 2. Data & State

- **Database:** PostgreSQL
  - *Why:* Robust, scalable, standard relational database. Easily hosted as a managed component on DigitalOcean App Platform.

- **Schema Management:** goose
  - *Why:* Manages database structure changes (migrations) using plain SQL files. Ensures production and dev DBs stay in sync.

- **Data Access:** sqlc
  - *Why:* Generates type-safe Go code from raw SQL queries. Replaces complex ORMs and catches SQL errors at compile time.

- **Real-time:** nhooyr.io/websocket + HTMX ws extension
  - *Why:* Robust, idiomatic Go WebSocket handling to push game state updates to clients instantly.

## 3. Infrastructure & Deployment

- **Host:** DigitalOcean App Platform
  - *Why:* Fully managed PaaS. Handles build, deployment, SSL, and scaling automatically. Zero server maintenance.

- **Containerization:** Dockerfile
  - *Why:* Standardizes the build environment for the App Platform.

- **Web Server / Proxy:** Managed by Platform
  - *Why:* Platform handles SSL termination (HTTPS) and routing automatically.

- **Deployment Pipeline:** Git Push -> Auto-Deploy
  - *Why:* Continuous Delivery. Pushing to the repository triggers a build and update.

## 4. Authentication (Strategy)

- **Phase 1 (MVP):** Anonymous Sessions
  - *Implementation:* Signed HTTP-only cookies using `alexedwards/scs`.
  - *Identity:* Users are identified by a random Session ID and a display name they pick per lobby.

- **Phase 2 (Future):** OAuth2 / Magic Links
  - *Implementation:* `markbates/goth` for OAuth.
  - *Why:* Avoids storing passwords (security liability). Accounts are linked to the existing Session IDs.

## 5. Development Tooling

- **Live Reload:** Air
  - *Why:* Recompiles and restarts the Go server instantly when files change.

- **Build Tool:** Make (Makefile)
  - *Why:* Universal standard to coordinate `templ generate`, `tailwind`, and `go build` steps.

## 6. Project Structure

```
mindmeld/
├── cmd/
│   └── server/
│       └── main.go           # Application entrypoint
├── internal/
│   ├── handlers/             # HTTP handlers
│   ├── middleware/           # Auth, logging, etc.
│   ├── db/                   # sqlc generated code
│   └── ws/                   # WebSocket hub
├── templates/                # templ files (.templ)
├── static/                   # CSS, JS, images
│   └── css/
│       └── input.css         # Tailwind source
├── migrations/               # goose SQL migrations
├── queries/                  # sqlc SQL queries
├── go.mod
├── go.sum
├── Makefile
├── Dockerfile
├── .air.toml
├── sqlc.yaml
└── tailwind.config.js
```

## 7. Key Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/go-chi/chi/v5` | HTTP router |
| `github.com/alexedwards/scs/v2` | Session management |
| `github.com/a-h/templ` | Type-safe HTML templates |
| `nhooyr.io/websocket` | WebSocket handling |
| `github.com/pressly/goose/v3` | Database migrations |
| `github.com/sqlc-dev/sqlc` | SQL code generation |
| `github.com/jackc/pgx/v5` | PostgreSQL driver |

## 8. Common Commands

```bash
# Development
make dev          # Start with live reload (Air)
make generate     # Run templ + sqlc code generation
make build        # Build production binary

# Database
make migrate      # Apply pending migrations
make migrate-new  # Create new migration file

# Styling
make css          # Build Tailwind CSS
make css-watch    # Watch mode for Tailwind
```

## 9. Environment Variables

| Variable | Description |
|----------|-------------|
| `DATABASE_URL` | PostgreSQL connection string |
| `SESSION_SECRET` | Secret key for cookie signing |
| `PORT` | Server port (default: 8080) |
| `ENV` | Environment: development, production |
