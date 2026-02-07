---
title: "Multi-Agent Local Environment Isolation"
status: planned
description: "Run many local Mindmeld instances in parallel with isolated state and Postgres-faithful behavior."
tags: [area/devex, area/backend, type/infra, tech/postgres]
priority: high
created: 2026-02-07
updated: 2026-02-07
effort: M
depends-on: []
---

# Multi-Agent Local Environment Isolation

## Problem / Intent

We want multiple concurrent local Mindmeld instances (for parallel AI agent threads) without database collisions. The local setup currently assumes one shared database and a single dev thread.

## Constraints

- Keep production-faithful PostgreSQL behavior for `pgx`, `sqlc`, and goose migrations.
- Avoid defaulting to one container per thread.
- Keep setup and cleanup fast for repeated agent runs.
- Keep the developer workflow simple and explicit.

## Proposed Approach

### Phase 1: Shared Postgres Container + Many Databases (Primary Path)

1. Keep one local Postgres container.
2. Create one database per thread/worktree (for example: `mindmeld_<suffix>`).
3. Standardize per-thread `DATABASE_URL` wiring.
4. Run migrations per isolated database.
5. Add helper commands for create/list/drop lifecycle.
6. Document the "start a new agent thread" workflow.

### Phase 2: PGlite Compatibility Spike (Secondary Path)

Time-box a spike to test whether PGlite is practical for this project.

Acceptance criteria:
- Goose migrations run end-to-end.
- `sqlc`/`pgx` queries behave as expected.
- Transaction and concurrency behavior is acceptable for app/runtime usage.
- Operational ergonomics are better than Phase 1.

Relevant docs signals to verify in-repo:
- `pglite-socket` can expose a Postgres-compatible socket for external clients.
- PGlite instances are single-connection, which may constrain runtime usage.

## Open Questions

- What thread/worktree naming convention should we enforce?
- Should lifecycle tooling live in `Makefile` targets or a dedicated `scripts/` tool?
- Do we auto-clean stale thread databases?
- Is PGlite viable for full local runtime, or only for narrow test flows?

## Decision Summary (FAQ)

- Is PGlite a good idea?
Yes, as a spike and optional profile. It is not yet the default path.

- Why is single-connection a concern?
Mindmeld runtime uses concurrent request handling and `pgxpool`; a single DB connection can serialize work within one server instance and diverge from production concurrency behavior.

- If each agent has its own PGlite instance, does that fix it?
It fixes cross-agent isolation. It does not remove per-instance serialization limits inside each running app.

- What should be the default local strategy now?
One shared Postgres container with one database per agent/worktree/thread.

- What about cloud development environments?
Use managed Postgres with one database per agent/worktree as the default. Keep PGlite as a speed-first fallback for lower-fidelity runs.

## Validation Plan

1. Run `make fmt` and `make lint`.
2. Run `go build -o bin/server ./cmd/server`.
3. Run E2E validation with `make e2e-screenshot` (dev server must already be running), then read generated PNGs in `e2e/screenshots/` to verify expected UI behavior after DB switching.
