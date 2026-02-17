---
title: "Work Party Prep (Launch Readiness)"
status: active
description: "Polish both games for a work social event: reduce friction, expand content, add juice."
tags: [area/product, type/polish]
priority: high
created: 2026-01-23
updated: 2026-02-12
effort: L
depends-on: []
---

# Work Party Prep (Launch Readiness)

**Goal:** Make Trivia and Cluster genuinely fun to play at a work social event. Polish over new mechanics.

## Phase 1: Stabilize (Done)
**Objective:** Eliminate visual glitches and create a smooth, app-like feel.

- [x] **Refactor WebSocket Refreshes**:
    - Changed `broadcastRoundAdvanced` and `broadcastQuestionRevealed` to HTMX trigger events.
    - Updated frontend `game_content` container to listen and `hx-get` refresh.

- [x] **Answer Status Polish**:
    - Smooth "Who has answered" updates via `hx-swap-oob="true"`.

## Phase 2: Trivia Question Friction
**Objective:** Remove "writer's block" so games start fast and questions are good.

- [ ] **Content Rating System** ([content-rating.md](./content-rating.md)):
    - Host selects audience (Adults / Work / Kids) at lobby creation.
    - All content filtered by rating. Needed before adding packs or AI generation.

- [ ] **Curated Question Packs**:
    - Pre-built themed decks (pop culture, science, history, "about each other", etc.).
    - Host picks a pack (or "player-written") as the question source.
    - Each pack tagged with content rating.
    - Solves friction immediately — no API cost, no latency, works offline.

- [ ] **AI Question Assist** (optional layer on top of packs):
    - "Generate Question" button on the submit form.
    - AI generates 1 question + answers, constrained by lobby content rating.
    - Player can edit before submitting. AI enhances, doesn't replace.
    - Platform-funded with rate limits for MVP.

## Phase 3: Cluster Content Overhaul
**Objective:** Make Cluster rounds provoke opinions and split the room.

- [ ] **Expand prompt-axis pool** (quantity):
    - Target: 30+ prompt-axis combinations (currently 3).
    - Enough for multiple full sessions without repeats.

- [ ] **Improve prompt quality** (spiciness):
    - Prompts should provoke genuine disagreement, not consensus.
    - Rubric: Would people argue about this after the reveal?
    - Tag each combination with content rating.

- [ ] **Add axis variety**:
    - More axis sets beyond the current 2.
    - Axes should create meaningful splits (not just agree/disagree).

## Phase 4: Cluster Reveal Juice
**Objective:** Make the centroid reveal a dramatic, memorable moment.

- [ ] **Reveal choreography**:
    - Animate points appearing on the plane.
    - Centroid drops in with emphasis.
    - Winner highlight with visual flourish.

- [ ] **Post-reveal insights**:
    - Show who was the outlier.
    - Optional discussion prompt ("Why did the group land here?").

## Phase 5: Small Polish
**Objective:** Quality-of-life improvements that add up.

- [x] **Trivia: show "unanswered" bar** in the revealed results distribution.
- [ ] **Game instructions** ([game-instructions.md](./game-instructions.md)): pre-game rules screen.
