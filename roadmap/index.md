# Mindmeld Roadmap

## North Star

**Promote conversation, don't capture attention.**

The game is a catalyst for human connection, not a substitute for it. Every feature should:
- Get people talking to each other, not staring at screens
- Create moments of laughter, debate, or surprise that happen *between* players
- Minimize time spent waiting, typing, or navigating UI
- End gracefully so people can move on to other activities

If a feature makes the game stickier but less social, it's wrong for Mindmeld.

---

## Goal

A platform for multiplayer party games that bring people together through shared thinking. Primary goal is social connection; secondary is stats tracking and visualization.

## Current Focus

**[Work Party Prep](./work-party-prep.md)** - Getting the platform polished and ready for a company social event.

---

## WSJF Prioritized Backlog

*Sorted by (Value + Urgency + Risk) / Size. Do these in order.*

| # | Item | Why High Priority |
|---|------|-------------------|
| 1 | **Refactor Flicker Fix** | Current `reload()` fix is jarring; needs HTMX/WS polish |
| 2 | **AI Question Assist** | Reduces friction for players submitting questions |
| 3 | **Round results screen** | Creates discussion moments between questions |
| 4 | **Show answer status while waiting** | Reduces "is it frozen?" anxiety |
| 5 | **Mid-game player join sync** | Prevents confusion when friends arrive late |
| 6 | **Per-round point analysis** | Fuels friendly banter about who's winning |
| 7 | **Visual redesign completion** | Polish remaining screens for cohesive feel |
| 8 | **Enhanced end-of-game stats** | Conversation fuel: "most missed question" etc |

---

## Active Work Units

- [work-party-prep.md](./work-party-prep.md) - **HIGHEST PRIORITY** - Stability & Polish
- [trivia.md](./trivia.md) - Trivia MVP implementation (core complete, polish remaining)
- [visual-redesign.md](./visual-redesign.md) - 1970s cerebral sci-fi UI theme
- [question-system.md](./question-system.md) - Question types, sources, AI generation, pools
- [housekeeping.md](./housekeeping.md) - Build cleanup & project structure
- [game-isolation-refactor.md](./game-isolation-refactor.md) - Refactor to isolate game logic (Trivia vs Platform)
- [security.md](./security.md) - App Security Hardening
- [analytics.md](./analytics.md) - Plausible analytics for gameplay + stats
- [codebase-map.md](./codebase-map.md) - Auto-generated docs & key touchpoints
- [chronology.md](./chronology.md) - Daily timeline sorting game (planned)
- [vector-golf.md](./vector-golf.md) - Vector-based semantic golf game (planned)
- [anonymous_match.md](./anonymous_match.md) - Anonymous Match social game (future)

## Key Links

- Production: (deployed via DO App Platform)
- [Project README](../README.md)
