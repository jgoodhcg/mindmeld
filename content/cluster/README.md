# Cluster Content Library

This directory is the source of truth for imported Cluster prompts/axes.

## Files

- `library.v1.json` - current content library

## Workflow

1. Edit `library.v1.json`.
2. Validate locally:
   - `go run ./cmd/cluster-content validate -file content/cluster/library.v1.json`
3. Preview DB changes (no writes):
   - `go run ./cmd/cluster-content import -file content/cluster/library.v1.json -dry-run -env dev`
4. Import to your dev DB when ready:
   - `go run ./cmd/cluster-content import -file content/cluster/library.v1.json -env dev`

## Dev vs Prod Imports

The importer supports both development and production databases with explicit safety checks.

- Dev import defaults:
  - `-env` defaults to `dev`.
  - local URLs like `localhost`, `127.0.0.1`, and `*.local` are treated as dev.
- Prod import requirements:
  - URL must look production-like (non-local host).
  - `-env=prod` must be set.
  - `-allow-production` must be set.

Example production preview:

- `go run ./cmd/cluster-content import -file content/cluster/library.v1.json -database-url "$DATABASE_URL_PROD" -env prod -allow-production -dry-run`

Example production import:

- `go run ./cmd/cluster-content import -file content/cluster/library.v1.json -database-url "$DATABASE_URL_PROD" -env prod -allow-production`

## Scaling to 500+ pairs

The model is prompt-centric: each prompt lists multiple `axis_slugs`.

- Pair count is roughly:
  - `sum(len(axis_slugs) for prompt in prompts)`
- To reach 500 quickly:
  - 125 prompts x 4 axis links each = 500 pairs
  - 100 prompts x 5 axis links each = 500 pairs

## Authoring guidance

- Keep prompts work-safe unless intentionally marked for higher ratings.
- Prefer prompts that create disagreement, not trivia-style factual answers.
- Avoid duplicates phrased slightly differently.
- Keep each prompt mapped to 3-6 axis sets for variety.

## Rating rules

- `min_rating` values:
  - `10` = Kids
  - `20` = Work
  - `30` = Adults

A pair is available only when both prompt and axis are allowed by lobby rating.
