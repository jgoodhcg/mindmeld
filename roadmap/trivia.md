# Trivia

Players submit questions for each other, answer them, and see how well they know each other and the world.

## MVP

The simplest playable trivia experience.

### Scope

**Included:**
- Create lobby, get shareable code/URL
- Join lobby with display name
- Each player submits 1 question (multiple choice, 4 options)
- Host starts game
- Questions shown one at a time
- Players answer (author cannot answer their own)
- After all players answer, show correct answer
- End screen: scoreboard showing X/Y correct per player

**Excluded from MVP:**
- Timers (auto-advance when all answer)
- Teams (solo only)
- Time-based scoring (just correct = 1 point, incorrect = 0)
- Negative points
- Complex stats
- Host configuration UI
- Real-time sync (use polling/refresh)
- Session reconnection

### Game Flow (MVP)

1. Host creates lobby → gets code
2. Players join via URL, enter name
3. Each player submits 1 question with 4 answers (mark correct one)
4. Host starts game
5. For each question:
   - Show question + answers to all
   - Players select answer (author blocked)
   - When all have answered → reveal correct answer
   - Next question
6. Show final scoreboard

### Data Model

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

## Future Enhancements

### Question Answering Modes

Different ways to answer questions, selectable by host:

| Mode | Description |
|------|-------------|
| **Basic** (MVP) | Select answer, no time pressure, 1 point for correct |
| **Timed** | Countdown timer per question, faster = more points |
| **Confidence** | Players wager points based on confidence before seeing options |

### Question Types

| Type | Description |
|------|-------------|
| **Multiple Choice** (MVP) | 4 options, 1 correct |
| **True/False** | Binary choice |
| **Free Text** | Type answer, fuzzy matching |
| **Numeric** | Closest guess wins |

### Teams

- Solo mode (MVP - everyone vs everyone)
- Auto-assign teams (random balanced)
- Manual team assignment by host

### Real-time Features

- WebSocket-based live updates
- See other players typing/answering
- Live scoreboard updates
- Reconnection handling

### Stats & Scoring

- Time-based point bonuses
- Negative points for wrong answers
- Per-question stats (avg response time, % correct)
- Per-player detailed stats
- Per-team aggregate stats

### Host Configuration

- Questions per person (1, 2, 3, etc.)
- Answer mode selection
- Timer duration (10s, 15s, 30s)
- Team mode selection
