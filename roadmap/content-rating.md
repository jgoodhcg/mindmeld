---
title: "Content Rating System"
status: ready
description: "Cross-cutting audience rating (Kids/Work/Adults) for filtering content across all games."
tags: [area/platform, type/feature]
priority: high
created: 2026-02-12
updated: 2026-02-12
effort: M
depends-on: []
---

# Content Rating System

## Intent

Let hosts set a content audience for their lobby so all surfaced content (trivia packs, AI-generated questions, cluster prompts/axes) is appropriate for who's actually in the room. The default is permissive; filtering is opt-in.

**North Star Check:** Reduces social friction ("is this appropriate?") so players focus on each other, not on policing content.

## Specification

### Rating Tiers

| Tier | Label | Description |
|------|-------|-------------|
| 3 | **Adults** | Default. Spicy takes, edgy humor, hot opinions, adult references. |
| 2 | **Work** | Professional-safe. No innuendo, controversial opinions kept mild, no alcohol/drug references. |
| 1 | **Kids** | Family-friendly. Simple language, no controversy, no mature themes. |

### Host UI

Lobby creation includes a "Who's playing?" selector:

```
Who's playing?
  * Adults        (default)
  * Work
  * Kids
```

Displayed as a simple radio group. Selection is stored on the lobby and cannot be changed mid-game.

### Content Tagging

All content sources get a `min_rating` field:

- `1` = appropriate for Kids and above
- `2` = appropriate for Work and above
- `3` = Adults only

The lobby's selected rating acts as a ceiling. Content is shown only when `content.min_rating <= lobby.rating`.

### Affected Systems

| System | How rating applies |
|--------|-------------------|
| **Trivia question packs** | Each pack (or individual question) tagged with `min_rating` |
| **AI question generation** | Rating passed as constraint in the generation prompt |
| **Cluster prompts + axes** | Each prompt-axis combination tagged with `min_rating` |
| **Future games** | Any content pool follows the same pattern |

### Data Model

```sql
-- Lobby-level rating
ALTER TABLE lobbies ADD COLUMN content_rating SMALLINT NOT NULL DEFAULT 3;
-- 1=Kids, 2=Work, 3=Adults

-- Trivia question packs (when packs are added)
-- Each pack row gets: min_rating SMALLINT NOT NULL DEFAULT 3

-- Cluster prompt-axis sets
ALTER TABLE coordinates_prompt_axis_sets ADD COLUMN min_rating SMALLINT NOT NULL DEFAULT 3;
```

## Validation

- [ ] Host can select rating during lobby creation
- [ ] Default is Adults (3) when no selection made
- [ ] Lobby rating persists and is visible to players
- [ ] Content queries filter by `min_rating <= lobby.rating`
- [ ] AI generation instructions vary by rating tier
- [ ] E2E: create lobby with Work rating, verify only Work-safe content appears

## Scope

- Lobby-level setting only (not per-player or per-round)
- No content moderation or reporting in this unit (future)
- No user-facing "rating" badge on individual content items
- Rating is set at lobby creation and immutable for the session

## Context

- Agreed framing: "who's in the room" not "how edgy is the content"
- Adults is the default because restriction is opt-in
- Needs to land before or alongside content expansion (packs, cluster prompts) so new content is tagged from the start
