---
title: "Cluster Content Studio"
status: draft
description: "File-first authoring and live review tool for Cluster prompts, with tabular sources that generate the import JSON."
created: 2026-02-21
updated: 2026-02-21
tags: [area/game, area/tooling, type/feature, tech/cli]
priority: high
---

# Cluster Content Studio

## Intent

Make Cluster content easy to author and review at 500+ prompts without building a full in-app admin system.

The current `content/cluster/library.v1.json` format is good for import/validation, but it becomes hard to skim, sort, and QA as the library grows. This work adds a file-first content tooling workflow that is human-reviewable and still deterministic.

## Specification

Add a matrix-reloaded-style local executable workflow under the existing `cluster-content` tool:

- `cluster-content build`
  - Reads tabular prompt source files (CSV/TSV) from `content/cluster/prompts/`
  - Generates `content/cluster/library.v1.json`
  - Preserves stable slugs and deterministic import behavior
- `cluster-content review`
  - Runs a local web UI for reviewing prompt content files
  - Supports live reload on file changes
  - Shows sortable/filterable table views (rating, theme, status, axis usage)
  - Provides coverage stats and warnings (duplicates, missing axis refs, thin rating buckets)
- `cluster-content validate`
  - Continues to validate the generated JSON library (existing behavior)
- `cluster-content import`
  - Continues to import generated JSON into dev/prod databases with safety guards (existing behavior)

### Authoring Source Format (proposed)

Keep axis sets in JSON for now and move prompts to row-based files:

- `content/cluster/prompts/mild.csv`
- `content/cluster/prompts/polite.csv`
- `content/cluster/prompts/adults.csv`

Proposed columns:

- `slug`
- `text`
- `min_rating`
- `axis_slugs` (pipe-separated)
- `theme`
- `status` (`draft|ready`)
- `notes` (optional)

### UX Goals (Review Tool)

- Fast scanning of hundreds of prompts
- Search by phrase or slug
- Filter by rating/theme/status
- Inspect axis distribution and coverage gaps
- Spot duplicates and near-duplicates before import
- Export or print a review-friendly snapshot if needed

## Validation

How to know this is done:

- [ ] `cluster-content build` generates `content/cluster/library.v1.json` deterministically from tabular sources
- [ ] `go run ./cmd/cluster-content validate -file content/cluster/library.v1.json` passes on generated output
- [ ] `cluster-content review` loads the source files and updates after file changes (live reload)
- [ ] Review UI supports filtering by rating/theme and text search
- [ ] Duplicate and missing-axis warnings are visible in the review UI
- [ ] Round-trip docs exist in `content/cluster/README.md` (author -> build -> validate -> dry-run import -> import)
- [ ] E2E validation step (if review UI is mounted inside the app): run `make e2e-screenshot ARGS=\"/cluster\"` after UI-related changes and verify screenshots in `e2e/screenshots/`

## Scope

Included:

- Local content authoring/review tooling for Cluster prompt libraries
- Tabular prompt source format and JSON generation pipeline
- Local review UI (file-first, no auth)
- Content QA warnings and coverage summaries

Not included:

- In-app authenticated admin UI for editing production content
- Multi-user editing, permissions, moderation workflows
- Direct DB editing from the review tool
- AI generation workflows (can be layered later)

## Context

- Current import pipeline exists and should remain the deploy path:
  - `cmd/cluster-content/main.go`
  - `internal/clustercontent/`
- Current source-of-truth file (machine-friendly, hard to skim at scale):
  - `content/cluster/library.v1.json`
- This work supports the broader Cluster content expansion goals tracked in:
  - `roadmap/cluster-improvements.md`

## Open Questions (draft only)

- Should prompts be split by rating (`mild/polite/adults`) or by theme with a single combined file plus `min_rating` column?
- CSV vs TSV vs Markdown table as the source format (TSV may be safer for punctuation-heavy prompt text)?
- Should the review UI be implemented inside `cluster-content` (single binary) or as a separate helper executable?
- Do we want a generated static HTML report artifact for async review in git/PRs?
