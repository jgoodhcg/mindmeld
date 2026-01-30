# Shared Agent Guidelines

This is the source of truth for all AI agents (Claude, Gemini, etc.) working on this project.

## Workflow Preferences
- **One Step at a Time:** Perform a single logical task, then stop and ask for validation/feedback. Do not chain multiple feature implementations or fixes together.
- **Validation First:** Always reflect on the current state and plan before executing.
- **User Control:**
    - **Never** run the dev server (`make dev`, `air`, `go run`) automatically.
    - **Never** run migrations (`make migrate`) automatically.
    - **Never** query the local database unless explicitly asked.
    - **Never** assume the state of the database; ask the user to verify if unsure.
- **Roadmap Driven:** Keep the roadmap up-to-date and reference it frequently.

## Development Context
- **Tooling:** Dev tools (air, templ, sqlc, goose) are pinned in `go.mod`. Always use `go run ...` or the `Makefile`, do not assume global installs.
- **Generated Files:** 
    - `templates/*_templ.go`
    - `internal/db/` (sqlc output)
    - `static/css/output.css`
    - These are **not committed**. Run `make generate` and `make css` to build them.

## Tech Stack
- **Server:** Go 1.25+ with `chi` router.
- **Frontend:** `templ` (type-safe HTML) + `HTMX` (coming soon) + `Tailwind CSS`.
- **Database:** PostgreSQL with `sqlc` for type-safe queries.
- **Migrations:** `goose`.

## Allowed Verification Commands
The agent **IS** permitted to run these commands to verify code compilation and generation:

| Command | Description |
|---------|-------------|
| `make generate` | Generates templ templates and sqlc DB code. Run this before building. |
| `make templ` | Generates only templ templates. |
| `make sqlc` | Generates only sqlc DB code. |
| `go build -o bin/server ./cmd/server` | Compiles the server binary. Fails if there are syntax/type errors. |
| `make css` | Builds Tailwind CSS (safe to run, though usually handled by user/watcher). |
| `make fmt` | Formats code using `go fmt` and `templ fmt`. Run this before finishing changes. |
| `make lint` | Runs `go vet` to catch common errors. |

## E2E Visual Validation (Playwright)
The agent **IS** permitted to run these commands to validate UI changes visually:

| Command | Description |
|---------|-------------|
| `make e2e-screenshot` | Take screenshot of homepage. Output: `e2e/screenshots/` |
| `make e2e-screenshot ARGS="/path"` | Screenshot a specific page (e.g., `/lobby/ABC123`). |
| `make e2e-flow` | Run create-lobby flow with screenshots at checkpoints. |
| `make e2e-flow ARGS="join CODE"` | Run join-lobby flow for a specific lobby. |
| `make e2e-flow ARGS="trivia CODE"` | Run trivia game flow for a specific lobby. |
| `make e2e-test` | Run Playwright smoke tests. |

**Usage:**
- **Run these commands proactively** after making UI changes to validate they work correctly.
- After running, **read the PNG files** in `e2e/screenshots/` to visually verify the UI state.
- If commands fail with connection errors (e.g., "Failed to load", "ECONNREFUSED"), the dev server is not running. **Ask the user to start it** with `make dev`.

**Prerequisites:**
- Server must be running (`make dev` started by user)
- First run: `make e2e-install` to install dependencies

## User-Only Commands
The agent must **NOT** run these unless explicitly instructed:

| Command | Description |
|---------|-------------|
| `make dev` | Starts the server with hot reload. |
| `make migrate` | Applies database migrations. |
| `make db-up` | Starts the Docker database. |