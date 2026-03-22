---
title: "Work Party Prep (Launch Readiness)"
status: active
description: "Polish both games for a work social event: reduce friction, expand content, add juice."
tags: [area/product, type/polish]
priority: high
created: 2026-01-23
updated: 2026-03-22
effort: L
depends-on: []
---

# Work Party Prep (Launch Readiness)

## Intent

Make Trivia and Cluster reliably fun at a work social event. Favor fast starts, legible state, safe-by-default content, and memorable reveal moments over new mechanics.

## Current Focus

Final game-day readiness pass after the major polish work already shipped.

Immediate execution order:

1. Trivia result integrity + reconnect clarity
2. AI assist handling for first-person / personal facts
3. Audience wording pass across Trivia and Cluster
4. Question pack collision decision
5. 6-8 player rehearsal and friction cleanup

## Completion Snapshot

Shipped and validated so far:

- [x] Stabilized HTMX/WebSocket refresh flow so Trivia and Cluster update without full-screen flicker.
- [x] Added smooth answer-status updates and reconnect-driven resync behavior.
- [x] Added disconnect grace handling plus automatic host fallback after grace expiry.
- [x] Shipped shared audience/content rating controls for Trivia and Cluster.
- [x] Shipped curated Trivia question packs filtered by audience.
- [x] Shipped AI trivia drafting with local fallback, mocked e2e coverage, busy state, and stronger prompt handling.
- [x] Expanded Cluster prompt-axis pool, improved prompt quality, and broadened axis variety.
- [x] Imported the expanded Cluster source library to production (`cluster-library-v1`: 103 prompts, 17 axes, 515 pairs) and deactivated the older `cluster-seed-v2` content.
- [x] Added Cluster reveal choreography plus post-reveal outlier/debrief insights.
- [x] Enabled 2-player Cluster rounds and added a reveal-time `Dist from you` comparison so quick demos and small-group play still feel informative.
- [x] Added unanswered-result bars, clickable home branding, intentional host handoff, baseline analytics, game instructions, and accessibility baseline fixes.

## Specification

### Slice 1: Trivia result integrity + reconnect clarity

- [x] Stop keying Trivia answers and results by answer text alone.
- [x] Use stable option identity for answer storage, correctness, and revealed result aggregation.
- [x] Make the reconnect grace window legible to hosts and players so "waiting on reconnect" is clearly distinct from a stalled game.
- [x] Re-check Trivia reconnect/resync behavior so returning players do not appear to rewind room state.

### Slice 2: AI assist quality follow-up

- [x] Mocked e2e path exists for AI generation without a live billable call.
- [x] Busy/loading state is visible while generation is running.
- [x] Keyboard flow reaches the generate button correctly.
- [x] AI draft and question-pack affordances are visually grouped.
- [x] Prompt guidance/examples are clearer.
- [x] Named-subject and familiarity prompts preserve subject identity instead of inventing facts.
- [x] Improve first-person fact handling so input like `my favorite fruit is blueberry` becomes a question using the player's name or a placeholder such as `[MY_NAME]`.
- [ ] Upgrade distractor generation so personal questions still produce plausible wrong answers.
- [ ] Allow iterative AI refinement of a draft instead of forcing full regeneration.
- [ ] Save authored or AI-assisted questions into a reusable personal bank.
- [ ] Resolve or explicitly defer the current live-provider blocker: OpenRouter timeout after auth with Gemini.

### Slice 3: Audience wording pass

- [ ] Replace shared `Polite mode` / `Prompt filter` copy with clearer audience-safety language in Trivia and Cluster.
- [ ] Keep one shared control if possible, since the same rating affects Cluster prompts plus Trivia packs and AI-generated Trivia.
- [ ] Verify the revised wording reads cleanly in create-lobby, waiting-lobby, and in-game surfaces.

### Slice 4: Question pack collision control

- [ ] Decide whether duplicate pack picks in the same round are acceptable or should be actively discouraged.
- [ ] If mitigation is needed, prefer a lightweight approach that reduces collisions without making pack selection feel restrictive.
- [ ] Candidate options:
  - Randomize pack ordering per player.
  - Show a rotating subset instead of the full pack.
  - Temporarily reserve or hide a template once selected.
  - Assign non-overlapping suggestions per round or player.

### Slice 5: Rehearsal pass and friction cleanup

- [ ] Run a final 6-8 player rehearsal across Trivia and Cluster.
- [ ] Capture every observed friction point, then ship only the fixes that materially improve game-day flow.
- [ ] Prefer clarity, recovery, and pacing fixes over new feature ideas.

### Follow-on (only if readiness work is complete)

- [ ] Execute the highest-value items from [juice-playbook.md](./juice-playbook.md) in small, testable slices without compromising game-day stability.

## Validation

- [ ] Run `make fmt` and `make lint` after each shipped slice.
- [ ] Run `go build -o bin/server ./cmd/server` after code changes.
- [ ] Run `make e2e-test` before closing a UI-affecting slice. The dev server must already be running; the user will start it.
- [ ] For Trivia submit/AI/audience-copy work, run `make e2e-flow ARGS="templates"` and review generated PNGs in `e2e/screenshots/`.
- [ ] For live Trivia flow/reconnect work, run `make e2e-flow ARGS="trivia CODE"` when a test lobby is available and review the screenshots.
- [ ] For Cluster-facing wording or rehearsal follow-ups, run `make e2e-flow ARGS="coordinates CODE"` when a test lobby is available and review the screenshots.

## Scope

- In scope: launch-readiness polish, clarity, recovery, pacing, content quality, and safe-by-default wording.
- Out of scope: new game modes, large database/content authoring systems, or risky pre-event architecture work.
- Out of scope: any fix that requires migrations, infra changes, or live-provider dependency unless the event-readiness benefit is clear and immediate.

## Context

- Canonical roadmap overview: [index.md](./index.md)
- Trivia polish backlog: [trivia.md](./trivia.md)
- Cluster MVP status: [cluster-mvp.md](./cluster-mvp.md)
- Player guidance shipped here: [game-instructions.md](./game-instructions.md)
- Optional flavor follow-ons: [juice-playbook.md](./juice-playbook.md)
- Prior audience/content-rating work: [archived/content-rating.md](./archived/content-rating.md)
