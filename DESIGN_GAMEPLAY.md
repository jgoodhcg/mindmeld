# Mindmeld - Game Experience Overview

## Vision

Mindmeld is a platform for multiplayer party games that bring people together through shared thinking. The games explore how aligned (or misaligned) players' minds are through trivia, word association, and other social challenges.

**First game: Trivia** - Players submit questions for each other, answer under time pressure, and see how well they know each other and the world.

**Future games:**
- Word association (Wavelength-style)
- Turn-based strategy (simple 2D)
- Other social/party games TBD

---

## Game 1: Trivia

### Core Concepts

#### Lobby
- Created by a host, generates a unique join code and shareable URL.
- Host configures questions per person and team mode (manual assignment, auto-assign, or solo).
- Lobby exists until explicitly cleaned up.

#### Players
- No global accounts in v1; players pick a display name per lobby.
- Each player belongs to a team (solo = team of one).
- Identity persists per device for reconnects/reloads.

#### Questions
- Format (v1): Multiple choice with 4 options.
- Future: Other question types (free text, true/false, etc.).
- Players submit questions during a submission phase.
- Submitters do not answer their own questions.

#### Scoring
- Correct answers earn points (faster = more points).
- Incorrect answers lose points (fixed penalty, TBD).
- Timer is fixed in v1 (configurable later).

#### Game Flow
1) Lobby created; host gets shareable URL.
2) Players join and enter display names; wait in lobby.
3) Host configures questions per person and team mode.
4) Teams formed (manual, auto, or solo).
5) Question submission phase for each player.
6) Host starts the game.
7) Rounds repeat for each question:
   - Question shown to all players.
   - Timer counts down.
   - Players answer (except the author).
   - Round ends on timer or all answers.
   - Show correct answer and current stats.
8) Game ends with final stats; players determine the winner.

#### Stats Displayed
Per team (and per player in detail):
- Total points; average and median points per question.
- Total, average, and median response times.
- Fastest and slowest responses.

---

## Roadmap (Experience-Focused)
- Foundation: create lobbies, join via code, pick names, form teams.
- Submission: players author questions with four options; track submission counts.
- Play: timed rounds, answer submission (author blocked), points displayed.
- Results: show per-round feedback and end-of-game stats; option to play again.
- Polish: mobile-friendly UI, reconnection handling, clear error and edge cases.
- Future platform: additional game types (word association, strategy), optional accounts, player history, pre-made or AI-generated question packs.

---

## Open Questions
- Timer duration (10s? 15s? 30s?).
- Point values (max for instant answer, penalty for wrong).
- Lobby expiration policy and max players.
- Max questions per person.
