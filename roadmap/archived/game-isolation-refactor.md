---
title: "Codebase Refactor: Game Isolation"
status: done
description: "Refactor to isolate game-specific logic and enable new games."
tags: [area/backend, type/refactor]
priority: high
created: 2026-01-19
updated: 2026-02-07
effort: L
depends-on: []
---

# Codebase Refactor: Game Isolation

**Goal:** Isolate "Trivia" specific logic into its own namespace/package to facilitate adding future game types (e.g., "Vector Golf", "Imposter").

## Intent

Decouple Trivia-specific handlers, models, and runtime behavior from generic server and lobby code so new game types can be added without increasing conditional complexity in shared paths.

## Specification

- Introduce an `internal/games/trivia` package boundary for trivia game logic.
- Keep core server responsibilities in `internal/server` (auth/session/lobby plumbing).
- Extract a game interface the server can delegate to for lifecycle/actions/state.
- Refactor routing/handler wiring so game-specific code is not concentrated in shared handlers.
- Preserve current Trivia behavior while isolating responsibilities.

## Validation

- [x] `make fmt`
- [x] `make lint`
- [x] `go build -o bin/server ./cmd/server`
- [x] Run targeted gameplay verification for Trivia flows after refactor (E2E if UI/flow changes)

Validation run on 2026-02-07:
- `make fmt` passed.
- `make lint` passed.
- `go build -o bin/server ./cmd/server` passed.
- `make e2e-test` passed (3/3 smoke tests).
- Targeted UI flow passed via `make e2e-flow ARGS="templates"` and screenshots reviewed in `e2e/screenshots/templates-00-homepage.png`, `e2e/screenshots/templates-01-create-form-filled.png`, `e2e/screenshots/templates-02-lobby-created.png`, and `e2e/screenshots/templates-03-templates-modal.png`.

## Scope

- Includes code organization and interface boundaries for existing Trivia logic.
- Excludes building a second game implementation in this work unit.
- Excludes unrelated feature work or gameplay changes.

## Context

- `internal/server/handlers_game.go`
- `internal/server/routes.go`
- `internal/server/handlers_ws.go`
- Trivia-related models/events currently under `internal/server` and generated DB access in `internal/db`

## Motivation
Currently, game logic (handlers, models, events) is intermingled in `internal/server` and the global `db` package. Adding a second game type would clutter `handlers_game.go` and potentially require complex switch statements or tight coupling.

## Proposed Changes

### 1. Package Structure
Refactor the `internal` directory to group game-specific logic:

```
internal/
├── games/
│   ├── trivia/
│   │   ├── engine.go       # Game loop logic
│   │   ├── handlers.go     # HTTP Handlers
│   │   ├── models.go       # Specific structs (if not in DB)
│   │   └── events.go       # Trivia specific payloads
│   └── shared/             # Shared game interfaces (Player, Lobby)
├── server/
│   └── ...                 # Core server (auth, lobby management)
```

### 2. Database Schema
*   Ensure table names are prefixed (done: `trivia_rounds`, `trivia_questions`).
*   Consider a JSONB `game_state` column in `lobbies` or `rounds` for flexible game-specific data if schemas diverge significantly.

### 3. Interface Extraction
Define a `Game` interface that the server uses to delegate requests:

```go
type Game interface {
    Start(lobby *db.Lobby) error
    HandleAction(action string, payload any) error
    GetState(playerID string) GameState
}
```

### 4. Route Namespace
Move routes to:
*   `/lobbies/{code}/trivia/...` -> `/lobbies/{code}/game/...` (Generic) OR
*   Keep specific routes but organize handlers better.

## Steps
1.  Audit `handlers_game.go` and identify pure Trivia logic.
2.  Create `internal/games/trivia`.
3.  Move handlers and helper functions.
4.  Refactor `routes.go` to mount game-specific sub-routers or handlers.
