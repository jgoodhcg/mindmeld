# Cluster Content Library

This directory contains the source files for imported Cluster prompts/axes.

## Files

- `library.v1.json` - generated import library (still stores axis sets; prompts are generated from `prompts/*.csv`)
- `prompts/` - row-based prompt authoring files (`.csv` or `.tsv`)

## Workflow

1. Edit prompt source files in `content/cluster/prompts/` (and axis sets in `library.v1.json` if needed).
2. Rebuild the generated import library:
   - `go run ./cmd/cluster-content build -prompts-dir content/cluster/prompts -file content/cluster/library.v1.json`
3. Validate locally:
   - `go run ./cmd/cluster-content validate -file content/cluster/library.v1.json`
4. Preview DB changes (no writes):
   - `go run ./cmd/cluster-content import -file content/cluster/library.v1.json -dry-run -env dev`
5. Import to your dev DB when ready:
   - `go run ./cmd/cluster-content import -file content/cluster/library.v1.json -env dev`

Notes:

- `build` reads prompt rows from `prompts/*.csv` or `prompts/*.tsv` and copies `version`, `created_by_label`, and `axis_sets` from the JSON template (`-axis-file`, defaulting to `-file`).
- Prompt rows with `status=draft` are excluded from the generated library.

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

- Keep prompts polite/work-safe unless intentionally marked for higher ratings.
- Prefer prompts that create disagreement, not trivia-style factual answers.
- Avoid duplicates phrased slightly differently.
- Keep each prompt mapped to 3-6 axis sets for variety.

## Prompt Source Columns

Supported prompt source columns:

- `slug`
- `text`
- `min_rating` (`10|20|30` or `mild|polite|adults`)
- `axis_slugs` (pipe-separated: `slug-a|slug-b`)
- `theme` (optional; review tooling metadata)
- `status` (`draft|ready`; defaults to `ready`)
- `notes` (optional; review tooling metadata)

## Rating rules

- `min_rating` values:
  - `10` = Mild (family-friendly)
  - `20` = Polite (work-safe / mixed company)
  - `30` = Adults

A pair is available only when both prompt and axis are allowed by lobby rating.
