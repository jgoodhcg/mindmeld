# Analytics System

## Work Unit Summary

**Status:** idea

**Problem/Intent:**
Add analytics to understand gameplay usage patterns and power player-facing stats.

**Constraints:**
- Keep write load small; prefer server-side instrumentation.
- Separate high-level analytics from game-specific, queryable stats.

**Proposed Approach:**
- **Plausible:** Basic page views and session tracking via standard JS tag inclusion only.
- **Postgres:** In-app analytics and game-specific stats:
  - Track monthly and daily active users (MAU/DAU) with game segmentation.
  - Track game-specific stats per day (rounds played, questions answered, etc.) for each game.
  - Define a minimal event schema (game_start, game_end, round_complete, question_answered).
  - Store events in Postgres with lightweight rollup queries or nightly aggregates.
  - Expose a simple internal dashboard and reuse aggregates for player stats.

**Open Questions:**
- Which metrics are critical for MVP (DAU, games played, average session length)?
- Do we need opt-in or account-level toggles?
- How do we partition data to keep queries fast over time?
- How to ensure each new game includes these analytics hooks from the start?

---

## Notes

Focus on plausible, non-invasive analytics that align with the project's social/connection goals.
