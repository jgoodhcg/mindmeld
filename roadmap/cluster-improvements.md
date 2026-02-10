---
title: "Cluster Improvements"
status: draft
description: "Post-MVP roadmap for content tooling, AI augmentation, and reveal/analysis enhancements."
tags: [area/game, type/roadmap]
priority: medium
created: 2026-02-10
updated: 2026-02-10
effort: L
depends-on: [cluster-mvp.md]
---

# Cluster Improvements

## Intent

Plan post-MVP improvements that increase discussion and bonding without pulling attention away from players.

## Product Guardrail

Use this test for every enhancement:

- Does this increase player-to-player conversation?
- Does this reduce confusion or increase delight during reveal?
- Does this avoid turning the experience into passive screen-watching?

If attention cost is high and discussion lift is low, defer it.

## Phase Plan (Do Not Implement Yet)

### Phase 1: Content Expansion System

- Bulk prompt-axis curation workflow and content QA process.
- Prompt quality rubric for discussion potential, ambiguity level, and social safety.
- Versioned seed packs for different contexts (work party, friends, mixed familiarity).

### Phase 2: User-Created Content

- User-added prompts/axis sets with lightweight validation and moderation controls.
- Lobby-level enable/disable policy for player-generated content.
- Provenance and authorship display for generated or user-authored content.

### Phase 3: AI-Assisted Content Generation

- AI-assisted prompt drafting and axis drafting workflows.
- AI-generated candidate prompts with approval queue (not auto-live by default).
- Prompt diversity controls (avoid repetitive framing).

### Phase 4: Reveal and Insight Enhancements

- Visual "juice" experiments:
  - optional hull/border around all submitted points
  - optional cluster region overlays
  - stronger reveal choreography/animation variants
- Insight experiments:
  - identify farthest outlier from centroid
  - AI-generated optional discussion starter about centroid placement
  - AI-generated optional cluster labels
  - session-level trend analysis summaries

## Candidate Features Moved from Prior Combined Roadmap

- Benchmark mode (AI or historical anchor points).
- 1-2 player support variants.
- Alternate scoring models (outlier bonus, weighted/hybrid scoring).
- Team mode and longitudinal trend analysis.
- Per-round free-text rationale paired with coordinate submission.

## Evaluation Matrix Template (Per Feature)

For each proposed enhancement, score before building:

- Discussion lift: low/medium/high
- Cognitive overhead: low/medium/high
- Screen-attention risk: low/medium/high
- Implementation risk: low/medium/high
- Rollout strategy: prototype/flagged/default

Only move forward when discussion lift clearly outweighs attention risk.

## Immediate Next Planning Deliverables

- Draft a concrete prompt-expansion execution plan after MVP close.
- Create a scoped implementation plan for Phase 2 and Phase 3 with dependencies, rollout flags, and measurement criteria.
