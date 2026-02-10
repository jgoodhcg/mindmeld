---
title: "Coordinates"
status: active
description: "Cluster MVP: real-time 2D social alignment game with centroid scoring and normalized content schema."
tags: [area/game, type/feature]
priority: high
created: 2026-01-11
updated: 2026-02-10
effort: M
depends-on: []
---

# Coordinates (Cluster)
Refer to the game in source code, UI, and git commits as "Cluster". "Coordinates" was a placeholder name for this game.

## Intent

Ship the second multiplayer game for Mindmeld: a 2D social alignment game where players plot the same prompt on a shared coordinate plane, then react to the reveal together.

This MVP optimizes for:
- Fast party loop with minimal setup
- Strong reveal moments (cluster vs disagreement)
- Clear scoring that is easy to explain
- A content model that preserves provenance for future user- and AI-authored content

## Specification

### Gameplay (MVP)

- Minimum players to start: `3`.
- No fixed game length. Host can continue rounds indefinitely until stopping or prompt pool exhaustion.
- Each round uses one prompt + one axis set pair.
- Round phases:
  1. Prompt and axis set shown on shared screen.
  2. Players submit one `(x, y)` coordinate from phone.
  3. Shared screen shows waiting state until all active players submit.
  4. All dots reveal simultaneously.
  5. Centroid is highlighted with a pulse animation.
  6. Round winner + cumulative standings shown, then host continues.
- When no unused active prompt/axis pair remains for the session, show a "You have played all available prompts" end state.

### UX Refinements (Current Iteration)

- Player submission uses an interactive coordinate plane (tap/click to plot) instead of numeric text inputs.
- Axis labels render directly on the coordinate plane for both submit and reveal states.
- Submit and reveal coordinate planes should share the same rendering component/style language.
- Marker semantics:
  - other players: neutral/gray
  - current player: cyan highlight
  - round winner(s): amber outline matching standings highlight
  - centroid: target-style marker with pulse
- Scoring explanation must be explicit in reveal UI:
  - round points = distance-to-centroid score (`100` near centroid, `0` farthest)
  - total points = cumulative sum across rounds

### Scoring (MVP)

- Centroid-only scoring (no outlier bonus, no benchmark point, no AI/global anchor).
- For each player:
  - `distance = sqrt((x - centroid_x)^2 + (y - centroid_y)^2)`
  - `max_distance = sqrt(2)` (unit square diagonal)
  - `round_points = round((1 - min(distance / max_distance, 1)) * 100)`
- Round winner is highest `round_points` (ties allowed).
- Game standings are cumulative sum of `round_points` across rounds.

### Content Model and Provenance (MVP)

- Use normalized prompt/axis design now (axis reuse is first-class).
- Seed a small development set with `3` active prompt-axis pairs.
- Seed content can be AI-generated, but must preserve who authored/seeded it and which model produced it.

### Data Model (MVP)

```sql
CREATE TABLE coordinates_axis_sets (
  id UUID PRIMARY KEY,
  x_min_label TEXT NOT NULL,
  x_max_label TEXT NOT NULL,
  y_min_label TEXT NOT NULL,
  y_max_label TEXT NOT NULL,
  created_by_kind TEXT NOT NULL CHECK (created_by_kind IN ('system', 'developer', 'user')),
  created_by_player_id UUID NULL REFERENCES lobby_players(id),
  created_by_label TEXT NULL,
  authoring_mode TEXT NOT NULL CHECK (authoring_mode IN ('manual', 'ai_assisted', 'ai_generated', 'imported')),
  generator_provider TEXT NULL,
  generator_model TEXT NULL,
  generator_prompt_version TEXT NULL,
  generator_run_id TEXT NULL,
  provenance JSONB NOT NULL DEFAULT '{}'::jsonb,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CHECK (
    (created_by_kind = 'user' AND created_by_player_id IS NOT NULL) OR
    (created_by_kind IN ('system', 'developer') AND created_by_player_id IS NULL)
  )
);

CREATE TABLE coordinates_prompts (
  id UUID PRIMARY KEY,
  prompt_text TEXT NOT NULL,
  created_by_kind TEXT NOT NULL CHECK (created_by_kind IN ('system', 'developer', 'user')),
  created_by_player_id UUID NULL REFERENCES lobby_players(id),
  created_by_label TEXT NULL,
  authoring_mode TEXT NOT NULL CHECK (authoring_mode IN ('manual', 'ai_assisted', 'ai_generated', 'imported')),
  generator_provider TEXT NULL,
  generator_model TEXT NULL,
  generator_prompt_version TEXT NULL,
  generator_run_id TEXT NULL,
  provenance JSONB NOT NULL DEFAULT '{}'::jsonb,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CHECK (
    (created_by_kind = 'user' AND created_by_player_id IS NOT NULL) OR
    (created_by_kind IN ('system', 'developer') AND created_by_player_id IS NULL)
  )
);

CREATE TABLE coordinates_prompt_axis_sets (
  id UUID PRIMARY KEY,
  prompt_id UUID NOT NULL REFERENCES coordinates_prompts(id) ON DELETE CASCADE,
  axis_set_id UUID NOT NULL REFERENCES coordinates_axis_sets(id) ON DELETE CASCADE,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (prompt_id, axis_set_id)
);

CREATE TABLE coordinates_rounds (
  id UUID PRIMARY KEY,
  lobby_id UUID NOT NULL REFERENCES lobbies(id) ON DELETE CASCADE,
  prompt_axis_set_id UUID NOT NULL REFERENCES coordinates_prompt_axis_sets(id),
  round_number INT NOT NULL,
  centroid_x DOUBLE PRECISION NULL CHECK (centroid_x >= 0 AND centroid_x <= 1),
  centroid_y DOUBLE PRECISION NULL CHECK (centroid_y >= 0 AND centroid_y <= 1),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (lobby_id, round_number)
);

CREATE TABLE coordinates_submissions (
  round_id UUID NOT NULL REFERENCES coordinates_rounds(id) ON DELETE CASCADE,
  player_id UUID NOT NULL REFERENCES lobby_players(id),
  x DOUBLE PRECISION NOT NULL CHECK (x >= 0 AND x <= 1),
  y DOUBLE PRECISION NOT NULL CHECK (y >= 0 AND y <= 1),
  submitted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (round_id, player_id)
);
```

### Prompt Selection Rules (MVP)

- Select from `coordinates_prompt_axis_sets` where all related records are active.
- Do not repeat a pair within a lobby session until all active pairs are exhausted.
- Host skip/reroll marks current pair as used for that session and selects next.

### Out of Scope (MVP)

- Benchmark mode (AI point or historical global average point)
- 1-2 player support modes
- User-facing prompt/axis authoring UI
- Alternate scoring systems (outlier bonus, weighted/hybrid scoring)
- Team mode and long-term trend analytics
- Per-round free-text answer/rationale inputs alongside plotted coordinates (candidate Phase 2 experiment)

## Validation

- [ ] `make fmt` and `make lint`
- [ ] `go build -o bin/server ./cmd/server`
- [ ] Add/update tests for:
  - centroid calculation
  - round point calculation
  - per-lobby no-repeat prompt selection
  - exhaustion end-state trigger
- [ ] E2E validation (dev server must already be running):
  - run `make e2e-flow ARGS="coordinates CODE"`
  - read generated PNGs in `e2e/screenshots/` and verify:
    - waiting state is visible while pending players remain
    - reveal animation shows simultaneous dots
    - centroid pulse is visually clear
    - round standings render correctly

## Scope

- Build complete playable Cluster MVP with host/player flows, persistence, scoring, and seeded content.
- Keep implementation aligned with existing trivia/lobby architecture and conventions.
- Do not include Phase 2 features in this work unit.

## Context

- Canonical roadmap: `roadmap/index.md`
- Existing game flow reference: `roadmap/trivia.md`
- Data access patterns: `queries/`, `internal/db/`, `migrations/`
- UI conventions: `templates/`, `static/css/`, `DESIGN_SYSTEM_GUIDE.md`
- Policy and validation guardrails: `AGENTS.md`, `AGENT_BLUEPRINT.md`
