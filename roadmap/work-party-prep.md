---
title: "Work Party Prep (Launch Readiness)"
status: active
description: "Polish stability and reduce friction for a work social event launch."
tags: [area/product, type/polish]
priority: high
created: 2026-01-23
updated: 2026-02-02
effort: L
depends-on: []
---

# Work Party Prep (Launch Readiness)

**Goal:** Polish the experience for a work social event. Prioritize stability and low friction over new game mechanics.

## Phase 1: Stabilize (Immediate)
**Objective:** Eliminate visual glitches ("flicker") and create a smooth, app-like feel.

- [x] **Refactor WebSocket Refreshes**:
    - Current state: `window.location.reload()` causes full page reload white flash.
    - Desired state: Use HTMX event triggers or OOB swaps to fetch new content without refreshing.
    - **Tasks**:
        - [x] Change `broadcastRoundAdvanced` to send an HTMX trigger event (e.g., `events.EventRoundAdvanced`).
        - [x] Change `broadcastQuestionRevealed` to send an HTMX trigger event.
        - [x] Update frontend `game_content` container to listen for these events and perform an `hx-get` to refresh the partial.

- [x] **Answer Status Polish**:
    - Ensure the "Who has answered" indicators update smoothly without layout shifts.
    - Added `hx-swap-oob="true"` to partial updates to ensure correct targeting.

## Phase 2: Enhance (Pre-Event)
**Objective:** Remove "writer's block" friction during the question submission phase.

- [ ] **AI Question Assist**:
    - Add a "Generate Question" button to the submit form.
    - **Backend**: Implement a simple LLM handler (using platform key for now) to generate 1 trivia question + answers.
    - **Frontend**: Button triggers `hx-post` to fetch values and populate the form inputs.
    - **Constraint**: Ensure questions are "Safe for Work" and generally accessible knowledge.

## Phase 3: Expansion (Post-Event)
**Objective:** Deeply engaging solo play with shareable scoring.

- [ ] **Chronology Game**:
    - Drag-and-drop timeline sorting.
    - Wikipedia data source.
    - Daily challenge mode.
