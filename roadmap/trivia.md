# Trivia MVP

## Work Unit Summary

**Status:** active

**Problem/Intent:**
Build the first playable game on the Mindmeld platform. Players submit questions for each other, answer them, and see how well they know each other and the world. This establishes the core lobby/player infrastructure that future games will reuse.

**Constraints:**
- No timers or time-based scoring (auto-advance when all answer)
- Solo play only (no teams)
- Simple scoring: correct = 1 point, incorrect = 0
- No real-time sync (use polling/refresh)
- No session reconnection handling
- No host configuration UI beyond starting the game

**Proposed Approach:**
Implement a minimal game loop: create lobby → join with name → submit questions → play through each question → show scoreboard. Keep the implementation simple and defer complexity (timers, teams, real-time) to future work units.

**Open Questions:**
- How should question order be determined when the game starts? (Random shuffle? Order by submission time?)
- What happens if a player disconnects mid-game?

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
lobbies
  - id (uuid, primary key)
  - code (varchar, unique)
  - host_player_id (uuid, foreign key)
  - phase (enum: waiting | submitting | playing | finished)
  - created_at (timestamptz)

players
  - id (uuid, primary key)
  - lobby_id (uuid, foreign key)
  - name (varchar)
  - created_at (timestamptz)

trivia_questions
  - id (uuid, primary key)
  - lobby_id (uuid, foreign key)
  - author_player_id (uuid, foreign key)
  - question_text (text)
  - correct_answer (varchar)
  - wrong_answer_1 (varchar)
  - wrong_answer_2 (varchar)
  - wrong_answer_3 (varchar)
  - display_order (integer, set when game starts)

trivia_answers
  - id (uuid, primary key)
  - question_id (uuid, foreign key)
  - player_id (uuid, foreign key)
  - selected_answer (varchar)
  - is_correct (boolean)
  - answered_at (timestamptz)
```

---

## Future Enhancements (separate work units when ready)

These are noted here for context but should become their own work units when prioritized:

**Question Answering Modes:** Timed mode with countdown, confidence wagering mode.

**Question Types:** True/false, free text with fuzzy matching, numeric (closest guess wins).

**Teams:** Random auto-assign, manual assignment by host.

**Real-time:** WebSocket live updates, see others typing/answering, reconnection handling.

**Stats & Scoring:** Time-based bonuses, negative points, detailed per-question and per-player stats.

**Host Configuration:** Questions per person, timer duration, mode selection.
