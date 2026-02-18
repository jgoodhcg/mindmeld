---
title: "Content Rating System"
status: active
description: "Cross-cutting audience rating via lookup tiers (Kids/Work/Adults) for filtering content across all games."
tags: [area/platform, type/feature]
priority: high
created: 2026-02-12
updated: 2026-02-18
effort: M
depends-on: []
---

# Content Rating System

## Intent

Let hosts set a content audience for their lobby so all surfaced content (trivia packs, AI-generated questions, cluster prompts/axes) is appropriate for who's actually in the room.

Implementation preferences finalized:
- Lookup-backed ratings with ordered numeric IDs for easy insertion later.
- Default lobby audience is **Work**.
- Host can adjust audience **before game start** (waiting phase only).
- Work tier should stay closer to corporate-safe.

**North Star Check:** Reduces social friction ("is this appropriate?") so players focus on each other, not on policing content.

## Specification

### Rating Tiers

| Tier | Label | Description |
|------|-------|-------------|
| 30 | **Adults** | Spicy takes, edgy humor, hot opinions, adult references. |
| 20 | **Work** | Default. Corporate-safe: mild topics only, avoid innuendo and mature references. |
| 10 | **Kids** | Family-friendly. Simple language, no controversy, no mature themes. |

### Host UI

Lobby creation includes a "Who's playing?" selector:

```
Who's playing?
  * Adults
  * Work          (default)
  * Kids
```

Displayed as a simple radio group.

Inside the lobby room header:
- All players see a quiet "Audience" label.
- Host sees a small control to change rating while lobby phase is `waiting`.
- Updates are blocked once game start has occurred (`phase != waiting`).

### Content Tagging

All reusable content sources get a `min_rating` field:

- `10` = appropriate for Kids and above
- `20` = appropriate for Work and above
- `30` = Adults only

The lobby's selected rating acts as a ceiling. Content is shown only when `content.min_rating <= lobby.rating`.

### Affected Systems

| System | How rating applies |
|--------|-------------------|
| **Trivia question packs** | Each pack (or individual question) tagged with `min_rating` |
| **AI question generation** | Rating passed as constraint in the generation prompt |
| **Cluster prompts + axes** | Prompts and axis sets each tagged with `min_rating`; pair is usable only if both pass |
| **Future games** | Any content pool follows the same pattern |

### Data Model

```sql
CREATE TABLE content_ratings (
  id SMALLINT PRIMARY KEY,          -- 10/20/30, ordered tiers
  code TEXT UNIQUE NOT NULL,        -- kids/work/adults
  label TEXT NOT NULL,
  description TEXT NOT NULL,
  is_default BOOLEAN NOT NULL
);

ALTER TABLE lobbies
  ADD COLUMN content_rating SMALLINT NOT NULL DEFAULT 20
  REFERENCES content_ratings(id);

ALTER TABLE coordinates_prompts
  ADD COLUMN min_rating SMALLINT NOT NULL DEFAULT 30
  REFERENCES content_ratings(id);

ALTER TABLE coordinates_axis_sets
  ADD COLUMN min_rating SMALLINT NOT NULL DEFAULT 30
  REFERENCES content_ratings(id);

ALTER TABLE trivia_questions
  ADD COLUMN min_rating SMALLINT NOT NULL DEFAULT 30
  REFERENCES content_ratings(id); -- inherited from lobby on author submit
```

## Validation

- [x] Host can select rating during lobby creation
- [x] Default is Work (20) when no selection made
- [x] Host can change rating while lobby phase is waiting
- [x] Host cannot change rating after game start
- [x] Lobby rating persists and is visible to players (quiet label)
- [x] Cluster selection filters by `prompt.min_rating <= lobby.rating` AND `axis.min_rating <= lobby.rating`
- [x] Trivia authored questions inherit `min_rating = lobby.content_rating` on submit
- [ ] AI generation instructions vary by rating tier (deferred until AI generation flow is implemented)
- [x] E2E: create lobby with Work rating, verify only Work-safe content appears

## Scope

- Lobby-level setting only (not per-player or per-round)
- No content moderation or reporting in this unit (future)
- No user-facing "rating" badge on individual content items
- No admin CRUD for rating tiers in this unit
- Tier IDs are spaced to allow future insertion (e.g., 15 between Kids and Work)

## Context

- Agreed framing: "who's in the room" not "how edgy is the content"
- Work is the default for coworker game-day readiness
- Lookup table + ordered numeric IDs keeps filtering simple while leaving room for future tiers
- Needs to land before or alongside content expansion (packs, cluster prompts) so new content is tagged from the start
