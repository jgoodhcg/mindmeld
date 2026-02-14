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

---

## Work Units

### Active
- [work-party-prep.md](./work-party-prep.md) - Launch readiness polish and stability
- [trivia.md](./trivia.md) - Trivia MVP (core complete, polish remaining)
- [cluster-mvp.md](./cluster-mvp.md) - Cluster MVP delivery and final completion checklist

### Ready
- [content-rating.md](./content-rating.md) - Cross-cutting audience rating (Kids/Work/Adults) for all content
- [game-instructions.md](./game-instructions.md) - Game rules and role-based guidance for players

### Draft
- [multi-agent-local-isolation.md](./multi-agent-local-isolation.md) - Parallel local instances with isolated Postgres state
- [security.md](./security.md) - App security hardening
- [cluster-improvements.md](./cluster-improvements.md) - Post-MVP Cluster roadmap (user content, AI, reveal enhancements)
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

---

## WSJF Prioritized Backlog

*Sorted by (Value + Urgency + Risk) / Size. Do these in order. See [work-party-prep.md](./work-party-prep.md) for the active execution plan.*

| # | Item | Status | Why |
|---|------|--------|-----|
| 1 | ~~Refactor Flicker Fix~~ | Done | HTMX/WS polish shipped |
| 2 | **Content Rating System** | Ready | Foundation for all content work; host picks Kids/Work/Adults |
| 3 | **Curated Trivia Packs** | Next | Kills question-writing friction without API dependency |
| 4 | **AI Question Assist** | Next | Optional layer: generates questions constrained by content rating |
| 5 | **Cluster Prompt Overhaul** | Next | 3 prompts is too few, and they're too safe — need 30+ with spice |
| 6 | **Cluster Reveal Juice** | Planned | Centroid reveal needs animation and drama |
| 7 | **Game instructions screen** | Ready | Pre-game rules for new players |
| 8 | **Accessibility pass** | Backlog | Contrast, keyboard use, screen reader flow |
| 9 | **Mid-game player join sync** | Backlog | Late arrivals get confused |
| 10 | **Enhanced end-of-game stats** | Backlog | Conversation fuel: "most missed question" etc |

---

## Quick Ideas

- Hive Mind (cooperative convergence)
- Polarity (spectrum game)
- Pecking Order (ranking/sorting)
- See [ideas.md](./ideas.md) for the full backlog

## Key Links

- Production: (deployed via DO App Platform)
- [Project README](../README.md)
