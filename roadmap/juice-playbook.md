---
title: "Game Juice Playbook (Not Boring, 1970s Sci-Fi)"
status: draft
description: "A concrete backlog of high-flavor, conversation-first polish ideas for Trivia and Cluster."
tags: [area/product, area/frontend, type/polish, type/design]
priority: medium
created: 2026-02-19
updated: 2026-02-19
depends-on: [visual-redesign.md, work-party-prep.md]
---

# Game Juice Playbook (Not Boring, 1970s Sci-Fi)

## Intent

Capture a high-signal backlog of "juice" features that make sessions memorable, surprising, and social without turning Mindmeld into attention sludge. These ideas are inspired by !Boring-style product energy but grounded in Mindmeld's north star and 1970s cerebral sci-fi visual language.

## Principles

- **Conversation over optimization:** juice should spark debate, not grind.
- **Meaningful drama:** suspense and reveal mechanics should feel earned.
- **Ritual over noise:** repeated game loops should feel ceremonial, not spammy.
- **Legible spectacle:** every visual effect must communicate game state.
- **Human-centered sci-fi:** warm analog-lab vibe, not neon casino UX.

## Specification

### Track A: Trivia Juice

- **Launch sequence intro card**
  - Before each round, show a short "ROUND TRANSMISSION" card with theme + one-line stakes.
- **Confidence signal**
  - Optional "how sure are you?" chip before answer submit (Not Sure / Pretty Sure / Locked).
  - Reveal confidence mismatch moments ("very confident, but missed").
- **Crowd miss spotlight**
  - If most players miss a question, show "CONSENSUS ERROR" callout.
- **Round host line**
  - After reveal, host gets optional one-tap prompts: "Defend your wrong answer", "Fast debrief", "Move on".
- **Score momentum bands**
  - Replace plain leaderboard deltas with "rising / stable / slipping" indicators.
- **Comeback detector**
  - If trailing player closes gap quickly, show subtle "COMEBACK SIGNAL" badge.
- **Final recap cards**
  - End of game includes 3 social cards:
    - "Most Divisive Question"
    - "Unexpected Expert"
    - "Closest Finish"
- **Question provenance labels**
  - Show tiny tag on reveal: `Player-written`, `Pack`, or `AI assist`.

### Track B: Cluster Juice

- **Reveal countdown ritual**
  - 3-2-1 text cadence before points/centroid render.
- **Centroid impact frame**
  - Brief pulse ring + analog "target lock" grid flash when centroid lands.
- **Narrative reveal captions**
  - Add dynamic line after each reveal:
    - "Tight consensus"
    - "Split room"
    - "One strong outlier"
- **Distance language upgrade**
  - Convert raw points into flavor labels:
    - "Locked In", "Near Orbit", "Far Drift"
- **Outlier confessional prompt**
  - If outlier exists, show "Outlier, explain your coordinates in one sentence."
- **Axis tension summary**
  - Short sentence post-reveal:
    - "Group leaned Practical over Aspirational, with moderate spread."
- **Round archetype tags**
  - Auto-tag round mood: `Consensus`, `Debate Fuel`, `Chaos`, `Dead Heat`.
- **Session arc recap**
  - End-of-session mini timeline of round centroids (small multiples).

### Track C: Cross-Game Social Juice

- **Pre-game lobby "mood calibrator"**
  - One-tap session tone: `Light`, `Spicy`, `Chaotic`.
  - Influences copy style and prompt selection where applicable.
- **Inter-round transmission snippets**
  - One-line atmospheric messages between rounds:
    - "Signal coherence rising."
    - "Crew disagreement detected."
- **Host deck tools**
  - Compact host control panel with explicit social options:
    - "Pause for discussion"
    - "Speed run next round"
    - "Run tiebreak question"
- **Session postcard**
  - At end, generate a shareable in-lobby summary card (no external posting required).
- **Crew identity layer**
  - Lightweight per-session titles:
    - "Navigator", "Wildcard", "Centroid Whisperer", "Fact Reactor"
  - No permanent profile gamification.

### Track D: 1970s Sci-Fi Flavor Layer

- **Typographic ceremony**
  - Use all-caps mono labels for mode changes (`REVEAL`, `RESULTS`, `TRANSMISSION`).
- **Analog motion vocabulary**
  - Favor easing and pulse styles that feel like lab instrumentation.
- **Warm caution accents**
  - Amber for irreversible actions and high-stakes moments.
- **Data-panel composition**
  - Keep information in structured "console panels", not floating playful widgets.
- **Audio-ready hooks (optional future)**
  - Reserve event hooks for subtle synth ticks (not required for MVP).

## Validation

- [ ] Run `make e2e-test` after each shipped juice slice.
- [ ] Run `make e2e-flow ARGS="coordinates NEW"` for Cluster juice slices.
- [ ] Run `make e2e-flow ARGS="templates"` for Trivia submit/assist UX slices.
- [ ] Review `e2e/screenshots/` to ensure style still matches the 1970s sci-fi direction.
- [ ] Verify each juice feature creates a conversation beat, not just visual decoration.

## Scope

- Not included: addictive loops, loot systems, or engagement dark patterns.
- Not included: persistent progression economy.
- Not included: major gameplay rewrites that remove the current core loops.
- Scope is polish and social amplification of existing game structures.

## Context

- North star and current priorities: `roadmap/index.md`
- Active execution roadmap: `roadmap/work-party-prep.md`
- Visual constraints: `roadmap/visual-redesign.md`
- Existing game instruction framing: `roadmap/game-instructions.md`

## Suggested Sequencing

1. **Cluster reveal upgrades** (countdown + captions + archetype tags)
2. **Trivia recap cards** (conversation-focused end-state)
3. **Cross-game transmission snippets** (shared flavor layer)
4. **Session postcard artifact** (host closes with a memorable summary)

## Open Questions

- How much copy variability is desirable before "theater" becomes repetitive?
- Should host controls for pacing be surfaced inline or behind an advanced toggle?
- Which juice ideas should be rating-aware for Kids/Work/Adults audiences?
