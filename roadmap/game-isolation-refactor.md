---
title: "Codebase Refactor: Game Isolation"
status: planned
description: "Refactor to isolate game-specific logic and enable new games."
tags: [area/backend, type/refactor]
priority: high
created: 2026-01-19
updated: 2026-02-02
effort: L
depends-on: []
---

# Codebase Refactor: Game Isolation

**Goal:** Isolate "Trivia" specific logic into its own namespace/package to facilitate adding future game types (e.g., "Vector Golf", "Imposter").

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
