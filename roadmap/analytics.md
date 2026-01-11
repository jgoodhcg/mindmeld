# Plausible Analytics

## Work Unit Summary

**Status:** idea

**Problem/Intent:**
Add lightweight, privacy-respecting analytics so we can understand gameplay usage patterns and power player-facing stats.

**Constraints:**
- No third-party trackers; keep data local to Postgres.
- Avoid PII; analytics should be aggregate and session-based.
- Keep write load small; prefer server-side instrumentation.

**Proposed Approach:**
- Define a minimal event schema (game_start, game_end, round_complete, question_answered).
- Store events in Postgres with lightweight rollup queries or nightly aggregates.
- Expose a simple internal dashboard and reuse aggregates for player stats.

**Open Questions:**
- Which metrics are critical for MVP (DAU, games played, average session length)?
- Do we need opt-in or account-level toggles?
- How do we partition data to keep queries fast over time?

---

## Notes

Focus on plausible, non-invasive analytics that align with the project's social/connection goals.
