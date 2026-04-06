---
title: "Mindmeld Roadmap"
goal: "Multiplayer party games that bring people together through shared thinking and conversation."
---

# Roadmap

## North Star

**Promote conversation, don't capture attention.**

The game is a catalyst for human connection, not a substitute for it. Every feature should:
- Get people talking to each other, not staring at screens
- Create moments of laughter, debate, or surprise that happen *between* players
- Minimize time spent waiting, typing, or navigating UI
- End gracefully so people can move on to other activities

If a feature makes the game stickier but less social, it's wrong for Mindmeld.

---

## Current Focus

- **[Work Party Prep](./work-party-prep.md)** - Getting the platform polished and ready for a company social event.
- **Immediate next execution:** Audience wording pass, Cluster prompt review/prune, question-pack collision decision, then a 6-8 player rehearsal and friction cleanup.

---

## Work Units

### Active
- [work-party-prep.md](./work-party-prep.md) - Launch readiness polish and stability
- [trivia.md](./trivia.md) - Trivia MVP (core complete, polish remaining)
- [cluster-mvp.md](./cluster-mvp.md) - Cluster MVP delivery and final completion checklist
- [game-instructions.md](./game-instructions.md) - Game rules and role-based guidance for players

### Draft
- [multi-agent-local-isolation.md](./multi-agent-local-isolation.md) - Parallel local instances with isolated Postgres state
- [security.md](./security.md) - App security hardening
- [cluster-improvements.md](./cluster-improvements.md) - Post-MVP Cluster roadmap (user content, AI, reveal enhancements)
- [cluster-content-studio.md](./cluster-content-studio.md) - File-first authoring + live review tool for large Cluster prompt libraries
- [juice-playbook.md](./juice-playbook.md) - Not Boring-style polish backlog for Trivia + Cluster within 1970s sci-fi constraints
- [chronology.md](./chronology.md) - Daily timeline sorting game
- [vector-golf.md](./vector-golf.md) - Semantic golf daily game

- [visual-redesign.md](./visual-redesign.md) - 1970s cerebral sci-fi UI theme
- [question-system.md](./question-system.md) - Question types, pools, AI generation
- [analytics.md](./analytics.md) - Gameplay analytics + stats
- [codebase-map.md](./codebase-map.md) - Auto-generated docs & touchpoints
- [anonymous_match.md](./anonymous_match.md) - Anonymous Match social game
- [infrastructure.md](./infrastructure.md) - Scaling for WebSockets and events
- [ideas.md](./ideas.md) - Raw concept backlog

### Archived
- [archived/housekeeping.md](./archived/housekeeping.md) - Build cleanup & project structure
- [archived/css-cache-busting.md](./archived/css-cache-busting.md) - CSS cache invalidation
- [archived/game-isolation-refactor.md](./archived/game-isolation-refactor.md) - Isolate game logic for future games
- [archived/content-rating.md](./archived/content-rating.md) - Cross-cutting audience rating (Kids/Work/Adults) for all content

---

## WSJF Prioritized Backlog

*Sorted by (Value + Urgency + Risk) / Size. Do these in order. See [work-party-prep.md](./work-party-prep.md) for the active execution plan.*

| # | Item | Status | Why |
|---|------|--------|-----|
| 1 | ~~Refactor Flicker Fix~~ | Done | HTMX/WS polish shipped |
| 2 | **Content Rating System** | Done | Foundation shipped: host audience control + filtering tiers in gameplay content selection |
| 3 | **Curated Trivia Packs** | Done | Pack-based, work-safe-first templates shipped for faster starts |
| 4 | **AI Question Assist** | Done | Generate button shipped with audience-safe constraints + local fallback |
| 5 | **Cluster Prompt Overhaul** | Done | Expanded Cluster content shipped and imported: 103 prompts, 17 axes, 515 prompt-axis pairs |
| 6 | **Cluster Reveal Juice** | Done | Staggered reveal motion + centroid drop + outlier/debrief insights shipped |
| 7 | **Game instructions screen** | Done | Role-based pre-game + in-game help shipped for both games |
| 8 | **Accessibility pass** | Backlog | Contrast, keyboard use, screen reader flow. Note: Safari tab order is not intuitive/consistent with other browsers |
| 9 | **Mid-game player join sync** | Backlog | Late arrivals get confused |
| 10 | **Enhanced end-of-game stats** | Backlog | Conversation fuel: "most missed question" etc |

---

## Quick Ideas

- Hive Mind (cooperative convergence)
- Polarity (spectrum game)
- Pecking Order (ranking/sorting)
- Buzzer (speed trivia / game show buzz-in)
- See [ideas.md](./ideas.md) for the full backlog

## Key Links

- Production: (deployed via DO App Platform)
- [Project README](../README.md)
