---
title: "Anonymous Match"
status: draft
description: "Social game where players match anonymous responses to authors."
tags: [area/game, type/feature]
priority: low
created: 2026-01-11
updated: 2026-02-07
effort: M
depends-on: []
---

# Anonymous Match

## Work Unit Summary

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
