# Mindmeld - Design Document

## Vision

Mindmeld is a platform for multiplayer party games that bring people together through shared thinking. The games explore how aligned (or misaligned) players' minds are - through trivia, word association, and other social challenges.

**First game: Trivia** - Players submit questions for each other, answer under time pressure, and see how well they know each other and the world.

**Future games:**
- Word association (Wavelength-style)
- Turn-based strategy (simple 2D)
- Other social/party games TBD

---

## Stack

- **Framework**: Next.js 15 (with custom server for WebSocket support)
- **Runtime**: Node.js
- **Database**: PostgreSQL
- **Hosting**: DigitalOcean App Platform
- **Real-time**: WebSockets via `ws` library on custom server

## Architecture Note

Next.js doesn't natively support WebSocket connections. We use a custom `server.js` that:
1. Creates a Node HTTP server
2. Attaches the WebSocket server (`ws`) for real-time game communication
3. Passes regular HTTP requests to Next.js handler

Single port serves both HTTP and WebSocket traffic.

---

## Game 1: Trivia

### Core Concepts

#### Lobby

- Created by a host, generates a unique join code
- Shareable URL contains the join code (e.g., `/trivia/ABC123`)
- Host configures:
  - Questions per person (how many each player submits)
  - Team mode: manual assignment, auto-assign, or solo (teams of 1)
- Lobby persists in database until explicitly cleaned up (cleanup strategy TBD)

#### Players

- No global authentication (v1) - just a display name per lobby
- Players join via URL, enter display name
- Each player belongs to a team (even in solo mode - team of 1)
- Identity persistence: server issues a signed, httpOnly lobby session token on join that encodes `lobby_id` + `player_id`; revalidate it on every HTTP/WS request to survive reloads/disconnects. Clients never choose their own IDs.

#### Questions

- **Format (v1)**: Multiple choice with 4 options
- **Future**: Other question types (free text, true/false, etc.)
- Players submit questions during the submission phase
- Submitters cannot answer their own questions

#### Scoring

- **Correct answer**: Positive points (faster = more points)
- **Incorrect answer**: Negative points (fixed penalty, TBD)
- Timer duration: Fixed value for v1 (configurable later)

#### Game Flow

1. **Lobby Created** - Host gets shareable URL
2. **Players Join** - Enter display name, wait in lobby
3. **Host Configures** - Sets questions per person, team mode
4. **Teams Formed** - Manual, auto, or solo
5. **Question Submission** - Each player submits their questions
6. **Game Starts** - Host triggers start
7. **Question Round** (repeats for each question):
   - Question displayed to all players simultaneously
   - Timer counts down
   - Players submit answers (except question author)
   - Round ends when timer expires or all eligible players answered
   - Show correct answer + current stats
8. **Game End** - Display final stats, players determine winner

#### Stats Displayed

Per team (and per player in detailed view):
- Total points
- Average points per question
- Median points per question
- Total response time
- Average response time
- Median response time
- Fastest question response
- Slowest question response

---

## Data Model (Draft)

Designed to support multiple game types. `game_type` field allows filtering/routing.

```sql
lobbies
  - id (primary key, uuid)
  - code (unique join code, varchar)
  - game_type (varchar, e.g., 'trivia', 'wavelength')
  - host_player_id (foreign key)
  - config (jsonb, game-specific settings)
  - status (enum: waiting | submitting | playing | finished)
  - created_at (timestamp)

players
  - id (uuid)
  - lobby_id (foreign key)
  - name (varchar)
  - team_id (nullable, foreign key)
  - created_at (timestamp)

teams
  - id (uuid)
  - lobby_id (foreign key)
  - name (varchar)

-- Trivia-specific tables

trivia_questions
  - id (uuid)
  - lobby_id (foreign key)
  - author_player_id (foreign key)
  - question_text (text)
  - correct_answer (varchar)
  - wrong_answer_1 (varchar)
  - wrong_answer_2 (varchar)
  - wrong_answer_3 (varchar)
  - display_order (integer, set when game starts)

trivia_answers
  - id (uuid)
  - question_id (foreign key)
  - player_id (foreign key)
  - selected_answer (varchar)
  - is_correct (boolean)
  - response_time_ms (integer)
  - points_awarded (integer)
  - answered_at (timestamp)
```

---

## Roadmap

### Platform Identity/Auth
- [ ] Lobby session token issuance/validation (HTTP + WebSocket handshake)
- [ ] `players.user_id` nullable column for future auth backfill
- [ ] Auth provider integration (Auth.js/Clerk/Supabase) for post-v1

### Phase 1: Foundation
- [ ] Initialize Next.js project with TypeScript
- [ ] Set up custom server with WebSocket support
- [ ] Configure PostgreSQL connection (local + production)
- [ ] Create database schema/migrations
- [ ] Verify local dev environment works
- [ ] Deploy to DO App Platform (hello world)

### Phase 2: Lobby System
- [ ] Create lobby API endpoint (generates code)
- [ ] Join lobby via URL (`/trivia/:code`)
- [ ] Player name entry
- [ ] Lobby state synchronization via WebSocket
- [ ] Host configuration UI (questions per person, team mode)

### Phase 3: Teams
- [ ] Solo mode (auto-create team of 1)
- [ ] Auto-assign teams
- [ ] Manual team assignment UI

### Phase 4: Question Submission
- [ ] Question submission form (question + 4 answers, mark correct)
- [ ] Track submission count per player
- [ ] Show submission progress to lobby

### Phase 5: Game Play
- [ ] Randomize question order
- [ ] Display question to all players
- [ ] Server-authoritative timer
- [ ] Answer submission (blocked for question author)
- [ ] Point calculation (time-based positive, fixed negative)
- [ ] Round results display

### Phase 6: Stats & End Game
- [ ] Calculate all stats (total, avg, median, fastest, slowest)
- [ ] Team aggregate stats
- [ ] Final results screen
- [ ] Option to play again with same lobby

### Phase 7: Polish
- [ ] Mobile-responsive CSS
- [ ] Reconnection handling
- [ ] Error states and edge cases
- [ ] Lobby cleanup (old/abandoned lobbies)

### Future: Platform Expansion
- [ ] User authentication (optional accounts)
- [ ] Game selection from homepage
- [ ] Second game type (word association)
- [ ] Player history/stats across games
- [ ] AI-generated questions
- [ ] Pre-made question packs

---

## Open Questions / TBD

- Exact timer duration (10s? 15s? 30s?)
- Point values (max points for instant answer, penalty for wrong)
- Lobby expiration policy
- Max players per lobby?
- Max questions per person limit?
- PostgreSQL provider for production (Supabase free tier, Neon free tier, or DO Managed)
- Auth path: Auth.js vs. hosted (Clerk/Supabase) when adding global accounts; token sharing with WS auth helpers
