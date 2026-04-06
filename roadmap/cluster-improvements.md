---
title: "Cluster Improvements"
status: draft
description: "Post-MVP roadmap for content tooling, AI augmentation, and reveal/analysis enhancements."
tags: [area/game, type/roadmap]
priority: medium
created: 2026-02-10
updated: 2026-04-06
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

## Recent Playtest Signals

From the 2026-03-30 playtest plus the 2026-04-06 reveal pass:

- Players understood the premise and naturally mapped the prompt onto the plane.
- The original rules ambiguity has now been addressed by shifting default Cluster to self-answer + compare, not predict-to-win.
- The reveal works better when centroid is treated as group context rather than the round's "target."
- There is still interest in making the game feel richer than pure center-distance readouts, but prompt quality remains the first bottleneck.
- Two-player rounds still need extra scrutiny because pair play benefits less from centroid-only interpretation.

Implication:
- First, improve prompt quality and validate the comparison-first reveal in longer sessions.
- Then evaluate deeper reveal layers that add conversation fuel without turning the game into analytics homework.

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
  - nearest-neighbor / "who thinks most like you on this axis?" callouts
  - cross-prompt similarity summaries ("you were consistently high-confidence / low-spice across rounds")
  - relative-to-other-games or player-profile comparisons, only if they clearly increase conversation more than screen attention

### Phase 5: Labeled Points & Cross-Session Visibility

**User point labels:**
- Players optionally attach a short text label to their coordinate submission.
- Labels appear on hover/tap during reveal, giving context for placement choices.
- Evaluation: Does labeling increase discussion ("Why did you put X there?") or add typing overhead?

**Cross-session point overlays:**
- Reveal phase can optionally show anonymized points from previous sessions for the same prompt-axis pair.
- Creates "social benchmark" context: "Here's where others placed this" without live players.
- Session-level toggle (off by default) to avoid confusing new players.
- Requires data model for historical submissions with opt-in sharing.

**Opt-in consent for anonymous cross-session sharing:**
- Players (or host) explicitly opt-in during game setup to allow their submissions to be shared anonymously with future sessions.
- Default is opt-out: no data leaves the session unless consented.
- Consent scope options:
  - Per-session: "Share our anonymous placements with future groups"
  - Per-player (account-level): "Always share my anonymous placements" (if accounts exist)
- Stored consent flag tied to submission records for audit/compliance.
- Clear UI explaining what's shared: coordinates only, no names, no labels unless separately consented.

**Discussion starter potential:**
- Compare your placement to historical consensus/outliers.
- "You're in the 90th percentile for uniqueness on this prompt."

## Game Modes

Alternate scoring/target modes that change the objective while keeping the same coordinate-plane interaction.

### Core Modes

**Centroid (current default)**
- Goal: Place your own answer on the plane, then compare it to the group's center.
- Reveal: Show group center, center distance, peer distance, and notable outliers without default winner framing.
- Discussion: "Why did the room land here?" / "Who saw this differently, and why?"

**Follow**
- Goal: Place closest to a designated target player each round.
- Target rotates each round (or random/host-picked).
- Scoring: Inverse distance to target's coordinate.
- Discussion: "Why did [target] place there?" + target explains their reasoning.

**Avoid**
- Goal: Place furthest from a designated target player.
- Target rotates each round (or random/host-picked).
- Scoring: Direct distance from target (further = more points).
- Discussion: "Where did you think [target] would go to avoid them?"

### Extended Modes

**Scatter**
- Goal: Maximize distance from all other players (anti-centroid).
- Scoring: Sum of distances to all other points, or distance from centroid.
- Discussion: "Who went rogue and why?"

**Chase**
- Goal: Place closest to a moving target that shifts mid-round.
- Target starts at one position, shifts to another after X seconds.
- Scoring: Distance to final target position.
- Discussion: "Did you chase the shift or hold your ground?"

**Mirror**
- Goal: Match a hidden AI-generated "ideal" point for the prompt.
- AI places a point based on semantic interpretation of the prompt.
- Scoring: Inverse distance to AI point.
- Discussion: "Do we agree with the AI's interpretation?"

**Pairs**
- Goal: Secret partners try to place close to each other.
- Each player is assigned a secret partner at game start.
- Scoring: Combined score from centroid + partner proximity.
- Discussion: "Did you find your partner?"

### Mode Selection UX

- Host selects mode at lobby creation (default: Centroid).
- Mode displayed on game screen and in lobby list.
- Some modes require minimum players (e.g., Pairs needs even count).
- Consider mode-specific prompt filtering (some prompts work better for certain modes).

### Evaluation Criteria

For each mode, assess:
- Discussion lift: Does the reveal spark conversation?
- Cognitive overhead: Is the objective clear without explanation?
- Social dynamics: Does it create fun tension or just confusion?
- Implementation: Data model changes, scoring logic, UI updates.

## Candidate Features Moved from Prior Combined Roadmap

- Benchmark mode (AI or historical anchor points).
- Solo-only variants, if single-player practice or content QA ever becomes useful.
- Team mode and longitudinal trend analysis.

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
