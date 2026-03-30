---
title: "Cluster MVP"
status: active
description: "Playable real-time Cluster game with shipped scope, validation baseline, and follow-up notes."
tags: [area/game, type/feature]
priority: high
created: 2026-01-11
updated: 2026-03-30
effort: M
depends-on: []
---

# Cluster MVP

## Intent

Ship a complete, conversation-first social alignment game where players place points on a shared coordinate plane and discuss the reveal.

## Completion Snapshot

Implemented and validated:

- [x] Cluster game type integrated with lobby lifecycle and game routing.
- [x] Database schema and migrations for prompts, axis sets, prompt-axis pairs, rounds, and submissions.
- [x] Seed content with normalized prompt/axis provenance fields.
- [x] Minimum 2-player start gate and host-controlled round flow.
- [x] Prompt/axis no-repeat selection per lobby session with randomized ordering among unused pairs.
- [x] Centroid-distance scoring (0-100), winner detection, and cumulative standings.
- [x] Session exhaustion handling when prompt-axis pool is consumed.
- [x] Interactive click/tap coordinate plane submission (no raw numeric UX).
- [x] Shared input/reveal plane rendering with on-plane axis labels.
- [x] Marker semantics: gray others, cyan self, amber winner outline, centroid target crosshair.
- [x] Clear scoring explanation and standings labels (`Round pts`, `Avg/round`, `Total pts`).
- [x] Reveal standings show viewer-relative `Dist from you` values so two-player rounds still have an immediate comparison metric.
- [x] Real-time updates without wiping in-progress selection while other players submit.
- [x] Audience-aware content filtering (Kids/Work/Adults tiers) wired into prompt-axis selection.
- [x] Deterministic game listing order on platform page.
- [x] Automated test coverage:
  - unit tests for centroid and scoring logic
  - cluster multiplayer e2e for 2-player lobby flow with distance readout and preserved local selection during live updates
  - e2e visual flow screenshots for cluster states/reveal

## MVP Completion Note

- [x] Expanded active Cluster prompt-axis pool beyond the original starter set.
  - Source library now contains 103 prompts, 17 axes, and 515 prompt-axis pairs.
  - Production import completed on 2026-03-22 under `cluster-library-v1`.
  - Older `cluster-seed-v2` production content was deactivated after import to avoid mixing pools during selection.

The original final MVP deliverable is complete. Keep this file as the historical implementation checklist until it is archived.

## Validation Baseline

- `make fmt`
- `make lint`
- `go build -o bin/server ./cmd/server`
- `go test ./...`
- `npm run test:cluster` (from `e2e/`)
- `make e2e-flow ARGS="coordinates NEW"` and screenshot review

## Context

- Delivery log: `roadmap/cluster-delivery-log.md`
- Roadmap index: `roadmap/index.md`
