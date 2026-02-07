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

### Planned
- [multi-agent-local-isolation.md](./multi-agent-local-isolation.md) - Parallel local instances with isolated Postgres state
- [game-isolation-refactor.md](./game-isolation-refactor.md) - Isolate game logic for future games
- [security.md](./security.md) - App security hardening
- [chronology.md](./chronology.md) - Daily timeline sorting game
- [vector-golf.md](./vector-golf.md) - Semantic golf daily game

### Ideas
- [visual-redesign.md](./visual-redesign.md) - 1970s cerebral sci-fi UI theme
- [question-system.md](./question-system.md) - Question types, pools, AI generation
- [analytics.md](./analytics.md) - Gameplay analytics + stats
- [codebase-map.md](./codebase-map.md) - Auto-generated docs & touchpoints
- [anonymous_match.md](./anonymous_match.md) - Anonymous Match social game
- [coordinates.md](./coordinates.md) - Coordinates social alignment game
- [infrastructure.md](./infrastructure.md) - Scaling for WebSockets and events
- [ideas.md](./ideas.md) - Raw concept backlog

### Archived
- [archived/housekeeping.md](./archived/housekeeping.md) - Build cleanup & project structure
- [archived/css-cache-busting.md](./archived/css-cache-busting.md) - CSS cache invalidation

---

## WSJF Prioritized Backlog

*Sorted by (Value + Urgency + Risk) / Size. Do these in order.*

| # | Item | Why High Priority |
|---|------|-------------------|
| 1 | **Refactor Flicker Fix** | Current `reload()` fix is jarring; needs HTMX/WS polish |
| 2 | **Game isolation refactor** | Enables adding new games without trivia coupling or handler sprawl |
| 3 | **AI Question Assist** | Reduces friction for players submitting questions |
| 4 | **Round results screen** | Creates discussion moments between questions |
| 5 | **Accessibility pass** | Improves contrast, keyboard use, and screen reader flow for everyone |
| 6 | **Show answer status while waiting** | Reduces "is it frozen?" anxiety |
| 7 | **Mid-game player join sync** | Prevents confusion when friends arrive late |
| 8 | **Per-round point analysis** | Fuels friendly banter about who's winning |
| 9 | **Visual redesign completion** | Polish remaining screens for cohesive feel |
| 10 | **Enhanced end-of-game stats** | Conversation fuel: "most missed question" etc |

---

## Quick Ideas

- Hive Mind (cooperative convergence)
- Polarity (spectrum game)
- Pecking Order (ranking/sorting)
- See [ideas.md](./ideas.md) for the full backlog

## Key Links

- Production: (deployed via DO App Platform)
- [Project README](../README.md)
