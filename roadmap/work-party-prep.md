---
title: "Work Party Prep (Launch Readiness)"
status: active
description: "Polish both games for a work social event: reduce friction, expand content, add juice."
tags: [area/product, type/polish]
priority: high
created: 2026-01-23
updated: 2026-03-16
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

- [x] **Disconnect handling baseline**:
    - WebSocket reconnects now trigger an automatic lobby/game resync.
    - Player list shows disconnected/reconnecting participants.
    - Trivia and Cluster stop waiting forever once a disconnected player exceeds the reconnect grace window.
    - Host role transfers automatically after the grace window if the current host never returns.

## Phase 2: Trivia Question Friction
**Objective:** Remove "writer's block" so games start fast and questions are good.

- [x] **Content Rating System** ([archived/content-rating.md](./archived/content-rating.md)):
    - Host selects audience (Adults / Work / Kids) at lobby creation.
    - Host can change audience while lobby is still waiting.
    - Lobby + content filtering implemented (`min_rating <= lobby.content_rating`) for Cluster and authored Trivia questions.

- [x] **Curated Question Packs**:
    - Pre-built themed decks (work essentials, product/tech, pop culture, quick brain boost, world snapshot).
    - Pack-first modal in question submission flow (players can still write their own).
    - Pack/template visibility filtered by lobby content rating.
    - Solves friction immediately — no API cost, no latency, works offline.

- [x] **AI Question Assist** (optional layer on top of packs):
    - "Generate Question" button on the submit form.
    - AI generates 1 question + answers, constrained by lobby content rating.
    - Player can edit before submitting. AI enhances, doesn't replace.
    - OpenRouter-enabled when explicitly configured; local fallback generator keeps flow working without network/API key.
    - Deterministic mocked e2e coverage added for AI assist so normal validation does not spend money or depend on live provider behavior.
    - Default OpenRouter model updated to `google/gemini-3.1-pro-preview`.
    - Prompt contract upgraded to handle generic topics, stated facts, and personal-question shells in one path.

- [ ] **AI assist quality + UX follow-up**:
    - [x] Add an e2e path that exercises LLM generation safely and deterministically (no accidental billable live call during normal validation).
    - [x] Add a visible busy state while AI draft generation is running.
    - [x] Fix keyboard flow so tabbing from the AI topic input reaches the generate button reliably.
    - [x] Visually pair "Question pack" and "AI draft" so players can immediately see one is curated/manual and the other is AI-assisted.
    - [x] Add clearer prompt-writing guidance and examples for the AI draft input.
    - Improve prompt handling for personal/familiarity prompts so named subjects are preserved (`Justin` stays `Justin`) and the system can generate question shells with placeholders instead of inventing fake facts.
    - Upgrade distractor generation so personal questions yield convincing alternatives when the user already knows the correct answer.
    - Allow iterative AI refinement on a drafted trivia question instead of forcing a full regenerate from scratch each time.
    - Save authored or AI-assisted trivia questions into a personal question bank so players can reuse and adapt them in future games.
    - Local env loading now uses shell-compatible `.env.local` semantics in `make` targets and startup, avoiding stale or malformed inherited provider credentials during local runs.
    - Current live-provider blocker after auth is OpenRouter timeout with Gemini.

## Phase 3: Cluster Content Overhaul
**Objective:** Make Cluster rounds provoke opinions and split the room.

**Current priority:** Final game-day rehearsal follow-ups (6-8 player dry run and friction cleanup).

- [x] **Expand prompt-axis pool** (quantity):
    - Target: 30+ prompt-axis combinations (migration seeds 60 combinations).
    - Enough for multiple full sessions without repeats.

- [x] **Improve prompt quality** (spiciness):
    - Prompts should provoke genuine disagreement, not consensus.
    - Rubric: Would people argue about this after the reveal?
    - Tag each combination with content rating.

- [x] **Add axis variety**:
    - More axis sets beyond the current 2.
    - Axes should create meaningful splits (not just agree/disagree).

## Phase 4: Cluster Reveal Juice
**Objective:** Make the centroid reveal a dramatic, memorable moment.

- [x] **Reveal choreography**:
    - Animate points appearing on the plane.
    - Centroid drops in with emphasis.
    - Winner highlight with visual flourish.

- [x] **Post-reveal insights**:
    - Show who was the outlier.
    - Optional discussion prompt ("Why did the group land here?").

## Phase 5: Small Polish
**Objective:** Quality-of-life improvements that add up.

- [x] **Trivia: show "unanswered" bar** in the revealed results distribution.
- [x] **Audience control styling pass**: radio options and lobby audience controls aligned with existing UI visual language.
- [x] **Baseline analytics wiring**: Plausible script loaded in shared layout for site-wide pageview capture.
- [x] **Game instructions** ([game-instructions.md](./game-instructions.md)): pre-game rules screen.
- [x] **Accessibility baseline pass**: added skip link/main landmark, AI assist live status messaging, and keyboard/screen-reader-friendly question-pack dialog behavior.

## Phase 6: Flavor Expansion (Planned)
**Objective:** Add conversation-first "juice" moments without drifting from the north star.

- [ ] Execute highest-value items from [juice-playbook.md](./juice-playbook.md) in small, testable slices.
