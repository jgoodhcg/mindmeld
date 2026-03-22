---
title: "Trivia MVP"
status: active
description: "Core trivia game MVP with real-time play and polish tasks."
tags: [area/game, type/feature, tech/websocket]
priority: medium
created: 2025-12-14
updated: 2026-03-22
effort: L
depends-on: []
---

# Trivia MVP

> **NOTE:** Immediate polish and stability tasks (flicker refactor, AI assist) are currently tracked in **[Work Party Prep](./work-party-prep.md)**. Please refer to that file for the active prioritized list.

## Work Unit Summary

**Problem/Intent:**
Build the first playable game on the Mindmeld platform. Players submit questions for each other, answer them, and see how well they know each other and the world. This establishes the core lobby/player infrastructure that future games will reuse.

**Constraints & Refinements:**
- No timers or time-based scoring (auto-advance when all answer)
- Solo play only (no teams)
- Simple scoring: correct = 1 point, incorrect = 0
- **Real-time:** Use HTMX WebSockets for live updates (no manual refreshing).
- **Shuffle:** Answers must be shuffled so "A" isn't always correct. Questions shuffled for playback.
- **Styling:** Basic Tailwind CSS polish.
- WebSocket reconnects resync lobby/game content; disconnected players get a grace window before they stop blocking progress
- No host configuration UI beyond starting the game

**Proposed Approach:**
Implement a minimal game loop: create lobby → join with name → submit questions → play through each question → show scoreboard. Keep the implementation simple and defer complexity (timers, teams, real-time) to future work units.

**Open Questions:**
- ~~How should question order be determined when the game starts?~~ → Random shuffle implemented
- ~~What happens if a player disconnects mid-game?~~ -> Mark disconnected immediately, allow a short reconnect window, then stop counting them toward blocking actions
- ~~How should host transfer work if host leaves?~~ -> After the disconnect grace window, transfer host to the earliest still-connected player

## TODO (Host & Players)

- [x] **Host transfer on disconnect**: If host leaves mid-game, automatically transfer host to the earliest still-connected player after the reconnect grace window
- [x] **Manual host transfer**: Allow host to transfer host role to another player via UI
  - Shared lobby capability used by both Trivia and Cluster.
  - First release supports waiting lobbies plus between-round / revealed states.
  - Choices are limited to currently connected players and require explicit confirmation.

---

## Completed

- [x] Database schema (lobbies, players, lobby_players, trivia_rounds, trivia_questions, trivia_answers)
- [x] Create lobby flow
- [x] Join lobby with name
- [x] Question submission phase
- [x] Game playback (show questions, collect answers, reveal correct)
- [x] Scoreboard display
- [x] WebSocket real-time updates
- [x] Answer shuffling
- [x] Basic Tailwind CSS styling
- [x] Host can start game
- [x] Audience rating inheritance for authored questions (`trivia_questions.min_rating` from lobby setting)
- [x] Curated question packs in template modal (pack-based grouping, work-safe-first ordering, content-rating filtering)
- [x] AI Question Assist button in submit flow (OpenAI when explicitly configured, local fallback otherwise)
- [x] Disconnect handling baseline: reconnects resync the lobby automatically, disconnected players surface in the player list, and stalled questions unblock after the grace window

---

## TODO (Polish & UX)

- [x] **Join by code on home page**: Add a simple code input field so players can join a lobby without needing a direct link
- [x] **Remove public lobby list**: Hide the lobby list to keep games private; show only a count of active lobbies (total + lobbies with players)
- [x] **Refine home stats (active)**: Define active lobbies based on active WebSocket sessions, not just lobbies with players
- [ ] **Refine home stats (total)**: Show total trivia games ever played (e.g., via trivia rounds), not total lobbies
- [x] **Host-only start button**: Only display the "Start Game" button to the host player
- [x] **Fix white background on scroll**: Body/html background color shows white when scrolling past content
- [ ] **Show answer status while waiting**: Display which players have answered in the players section (for both question author and players who have already answered)
- [x] **Fix duplicate answer identity**: Stop keying submitted answers and result bars by answer text alone
  - Stable option identity (slot / answer ID) now drives answer submission, correctness checks, and results aggregation.
  - Result bars expose per-slot identity so duplicated answer text no longer causes mirrored vote counts.
  - Added e2e coverage with intentionally duplicated option text so this does not regress.
- [ ] ~~**Fix answer flicker**: UI flickers when other players submit answers (likely WebSocket update causing full re-render)~~ (Being refactored in [Work Party Prep](./work-party-prep.md))
- [x] **Handle ties on scoreboard**: Display a tie instead of arbitrarily choosing a winner when scores are equal
- [ ] **Better tie visualization**: Visually distinct handling for tied round winners (multiple crowns) and tied game winners
- [ ] **Author live results**: Allow the question author to see the live answer distribution graph while others are answering, instead of the static "Your Question" screen.
- [x] **Disconnect grace UX clarity**: Make the reconnect window legible to host and players
  - [x] During the grace window, show that a disconnected player is still temporarily blocking progress and can still rejoin.
  - [x] Add explicit host/player-facing copy for the waiting state, including the blocking player name and grace-window messaging.
  - [x] Verify reconnect behavior does not appear to "rewind" the question flow when a player returns during grace.

- [ ] **Personal fact placeholder handling in AI assist**: Improve prompt interpretation for first-person fact inputs
  - [x] If a player types a fact like `my favorite fruit is blueberry`, generate a question that either preserves the player's name or uses a safe placeholder like `[MY_NAME]`.
  - Avoid turning first-person facts into awkward generic trivia phrasing or inventing a third-person subject.
  - Cover both named-subject prompts and unnamed first-person prompts in mocked AI assist tests.

## TODO (Game Flow)

- [ ] **Round results screen**: After each question, show a breakdown of who answered what before moving to the next question (adds a "results" phase between questions)
- [ ] **Per-round point analysis**: Show point breakdown per question/round - who got what right, running totals, etc.
- [ ] **Enhanced end-of-game stats**: More detailed statistics at the end (most missed question, player accuracy %, question difficulty rankings, etc.)
- [ ] **Mid-game player join sync**: When a new player joins mid-game, broadcast update to all connected clients so player list stays in sync
- [ ] **Consider batch answering mode**: Alternative flow where players answer all questions at once, with progress shown, then reveal results together (vs current one-at-a-time with enforced sync)

---

## Game Flow

1. Host creates lobby → gets code
2. Players join via URL, enter name
3. Each player submits 1 question with 4 answers (mark correct one)
4. Host starts game
5. For each question: show to all → players select answer (author blocked) → when all answered → reveal correct → next
6. Show final scoreboard

## Data Model

```sql
lobbies (existing table, extended)
  + game_type (varchar, default 'trivia')
  + phase (varchar: waiting | playing | finished)

players
  - id (uuid, primary key)
  - device_token (varchar, unique) -- survives refreshes, enables rejoin
  - user_id (uuid, nullable)       -- future auth integration
  - created_at (timestamptz)

lobby_players
  - id (uuid, primary key)
  - lobby_id (uuid, fk lobbies)
  - player_id (uuid, fk players)
  - nickname (varchar)
  - is_host (boolean)
  - joined_at (timestamptz)
  - unique(lobby_id, player_id)
  - unique(lobby_id, nickname)

trivia_rounds
  - id (uuid, primary key)
  - lobby_id (uuid, fk lobbies)
  - round_number (integer)
  - phase (varchar: submitting | playing | finished)
  - created_at (timestamptz)
  - unique(lobby_id, round_number)

trivia_questions
  - id (uuid, primary key)
  - round_id (uuid, fk trivia_rounds)
  - author (uuid, fk players)
  - question_text (text)
  - correct_answer, wrong_answer_1/2/3 (varchar)
  - display_order (integer, set when round starts)
  - created_at (timestamptz)

trivia_answers
  - id (uuid, primary key)
  - question_id (uuid, fk trivia_questions)
  - player_id (uuid, fk players)
  - selected_answer (varchar)
  - is_correct (boolean)
  - answered_at (timestamptz)
  - unique(question_id, player_id)
```

**Design notes:**
- Device-token identity allows players to rejoin after refresh without re-entering name
- Rounds abstraction enables "play again" without creating a new lobby
- Host tracked via `is_host` flag on lobby_players (easy to transfer, guaranteed valid participant)
- Two-level phase: lobby phase (waiting/playing/finished) and round phase (submitting/playing/finished)

---

## Future Enhancements (separate work units when ready)

These are noted here for context but should become their own work units when prioritized:

**Question Answering Modes:** Timed mode with countdown, confidence wagering mode.

**Question Types:** True/false, free text with fuzzy matching, numeric (closest guess wins).

**Teams:** Random auto-assign, manual assignment by host.

**Spectator Mode:** Allow users to join a lobby as non-players to watch the game progress without participating.

**Real-time (Advanced):** Live typing indicators, "User X answered" notifications. Partial DOM updates for answer progress (avoid full page refresh that could interrupt players mid-answer).

**Stats & Scoring:**
- **Round Summary:** Breakdown of who picked what after each question.
- **Game Recap:** Fun visualizations (e.g., "Most confusing question", "Speediest player").
- Time-based bonuses, negative points.

**Host Configuration:** Questions per person, timer duration, mode selection.
