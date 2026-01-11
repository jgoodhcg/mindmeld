# Anonymous Match

## Work Unit Summary

**Status:** idea

**Problem/Intent:**
Build a social game that creates conversation by having players recognize each other through anonymous written responses. Players respond privately to a shared prompt, see responses anonymously, then try to match each response to the author. The fun comes from discussion, misattribution, and reveal—not speed or cleverness.

**Constraints:**
- Phones used only for submission and matching
- Conversation must happen off-screen
- AI must be invisible and non-authoritative
- Game works best with 4–7 players
- Optimize for laughter, recognition, and discussion—not optimization or strategy
- Facilitator-first game, not competitive system

**Proposed Approach:**
Implement a moderated prompt system, private response submission with AI validation for basic quality control, timer-driven conversation phase, anonymous matching phase, and dramatic reveal phase with one-by-one response revelation.

**Open Questions:**
- What is the approved prompt set source? (AI-generated vs curated vs community?)
- How should prompt selection work? (Random? Host choice?)
- What happens if a player's response is rejected by AI validation?
- Should there be a minimum number of players required?
- What should the conversation phase timer duration be?

---

## Game Flow

1. **Prompt Selection**
   - System (AI) selects one prompt from an approved prompt set
   - Prompt types: personal preference, opinion-with-reason, vibe translation, mild hypothetical
   - Prompt must be safe, open-ended, non-trivia

2. **Private Response Phase**
   - Each player submits a short text response (1–2 sentences)
   - Submissions must be: intelligible, on-topic, not substantially similar to another response
   - AI rejects nonsense, duplicates, or off-topic responses
   - No scoring or feedback at this stage

3. **Conversation Phase**
   - No devices
   - Players talk freely about the prompt
   - Players may express opinions or stories but may not quote their written response verbatim
   - This phase has a timer but no structured turns

4. **Matching Phase**
   - Screen displays all responses anonymously
   - Each player privately assigns each response to a player using their device
   - Timer runs during matching
   - All assignments lock at timer end

5. **Reveal Phase**
   - Responses are revealed **one at a time**:
     - Show the response
     - Reveal the author
     - Show which players guessed correctly
     - Pause briefly for reaction/conversation
   - Continue until all responses are revealed

6. **End Screen**
   - Show total correct guesses per player
   - Optionally show fun stats (e.g., "most misattributed")
   - No long-term progression required

---

## Scoring Rules (Optional but Recommended)

- +1 point for each correct match
- Optional: response author gains +1 if more than half the group guessed incorrectly
- No penalties for wrong guesses

---

## Data Model

```sql
lobbies (extended, reuse trivia structure)
  + game_type (varchar: 'trivia' | 'anonymous_match')
  + phase (varchar: waiting | prompt_shown | responding | conversing | matching | revealing | finished)

players (reuse from trivia)

lobby_players (reuse from trivia)

anonymous_match_prompts
  - id (uuid, primary key)
  - prompt_text (text)
  - prompt_type (varchar: preference | opinion | vibe | hypothetical)
  - is_approved (boolean)
  - created_at (timestamptz)

anonymous_match_rounds
  - id (uuid, primary key)
  - lobby_id (uuid, fk lobbies)
  - prompt_id (uuid, fk anonymous_match_prompts)
  - response_deadline (timestamptz)
  - conversation_end_time (timestamptz)
  - matching_deadline (timestamptz)
  - created_at (timestamptz)

anonymous_match_responses
  - id (uuid, primary key)
  - round_id (uuid, fk anonymous_match_rounds)
  - author_id (uuid, fk players)
  - response_text (text)
  - is_approved (boolean) -- AI validation result
  - rejection_reason (varchar, nullable) -- optional feedback for rejected responses
  - display_order (integer)
  - created_at (timestamptz)

anonymous_match_guesses
  - id (uuid, primary key)
  - round_id (uuid, fk anonymous_match_rounds)
  - guesser_id (uuid, fk players)
  - response_id (uuid, fk anonymous_match_responses)
  - assigned_player_id (uuid, fk players)
  - is_correct (boolean)
  - created_at (timestamptz)
  - unique(round_id, guesser_id, response_id)
```

**Design notes:**
- Reuse lobby/player infrastructure from trivia
- Prompt table supports future expansion and curation
- Response approval flag allows for AI moderation workflow
- Display order for randomized anonymous presentation
- One-to-one mapping between responses and guesses (each player guesses each response)
- Timer fields stored for client-side countdown display

---

## What This Is Not

- Not trivia
- Not deception-heavy
- Not real-time (except phase transitions)
- Not clever-wordplay-focused
- Not content-moderation-heavy

---

## Future Enhancements (separate work units when ready)

These are noted here for context but should become their own work units when prioritized:

**Prompt Management:** Host choice from pool, difficulty ratings, prompt categories, community submission and voting.

**Advanced AI Validation:** Contextual analysis, humor detection, more sophisticated duplicate detection.

**Game Variants:** 
- "Double Blind" - no conversation phase, straight to matching
- "Story Time" - extended responses (3-5 sentences) with storytelling focus
- "Hot Takes" - controversial prompts with debate-focused conversation phase

**Stats & Analytics:**
- "Most predictable player" (highest correct-guess rate)
- "Most mysterious player" (most misattributed)
- Prompt performance metrics (which prompts generate best discussions)
- Response length vs matching difficulty correlation

**Social Features:**
- "Highlight reel" of memorable responses after game
- Shareable summary card with funniest moments
- Favorite prompts saved per group

**Session Improvements:**
- Play again with same lobby without rejoining
- Multi-round sessions (3 prompts per session)
- Configurable timers per phase
