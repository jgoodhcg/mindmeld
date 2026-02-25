# Cluster Content Library

This directory contains the source files and generated artifacts for Cluster content imports.

## Files

- `source/` - canonical editable source files (`meta.json`, `axes.tsv`, `prompts.tsv`)
- `library.v1.json` - generated/import-compatible library snapshot (optional build artifact)

## Workflow

1. Edit source files:
   - prompts: `content/cluster/source/prompts.tsv`
   - axes: `content/cluster/source/axes.tsv`
   - metadata: `content/cluster/source/meta.json`
2. Validate directly from source (no JSON required):
   - `go run ./cmd/cluster-content validate -source-dir content/cluster/source`
3. Preview DB changes directly from source (no writes):
   - `go run ./cmd/cluster-content import -source-dir content/cluster/source -dry-run -env dev`
4. Import to your dev DB when ready:
   - `go run ./cmd/cluster-content import -source-dir content/cluster/source -env dev`
5. Optional: generate/update JSON snapshot artifact:
   - `go run ./cmd/cluster-content build -source-dir content/cluster/source -file content/cluster/library.v1.json`

Notes:

- `build`, `validate`, and `import` all support `-source-dir` (canonical TSV source pipeline).
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

## Source Files

### `source/meta.json`

- `version`
- `created_by_label`

### `source/axes.tsv`

Columns:

- `slug`
- `x_min_label`
- `x_max_label`
- `y_min_label`
- `y_max_label`
- `min_rating`

### `source/prompts.tsv`

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
