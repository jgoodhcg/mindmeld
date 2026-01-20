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

Building the **Trivia MVP** - the first playable game on the platform.

---

## WSJF Prioritized Backlog

*Sorted by (Value + Urgency + Risk) / Size. Do these in order.*

| # | Item | Why High Priority |
|---|------|-------------------|
| 1 | ~~**Fix answer flicker**~~ | ✅ Done |
| 2 | **Round results screen** | Creates discussion moments between questions |
| 3 | **Show answer status while waiting** | Reduces "is it frozen?" anxiety |
| 4 | **Mid-game player join sync** | Prevents confusion when friends arrive late |
| 5 | **Per-round point analysis** | Fuels friendly banter about who's winning |
| 6 | **Visual redesign completion** | Polish remaining screens for cohesive feel |
| 7 | **Enhanced end-of-game stats** | Conversation fuel: "most missed question" etc |
| 8 | **Refine home stats** | Low urgency, nice-to-have vanity metrics |
| 9 | **Batch answering mode** | Big change, needs validation before building |
| 10 | **Security hardening** | Important but not blocking current usage |
| 11 | **Codebase documentation** | Developer experience, not user-facing |
| 12 | **Analytics** | Useful but not urgent for small user base |

---

## Active Work Units

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
