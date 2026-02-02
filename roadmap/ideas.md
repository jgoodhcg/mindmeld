---
title: "Idea Backlog"
status: idea
description: "Raw game concepts awaiting promotion to work units."
tags: [area/product, type/idea]
priority: low
created: 2026-01-11
updated: 2026-02-02
effort: XS
depends-on: []
---

# Ruminating Idea List

A holding pen for game concepts. These are not fully specced work units yet.

---

## Social / Party Games (Sync)
*Focus: Connection, laughter, "Screen as Host"*

### 1. Hive Mind (Cooperative Convergence)
*   **Concept:** Can the group think exactly alike?
*   **Loop:** Category given -> everyone types 1 word -> reveal. If match: win. If not: round 2 (no talking), try to converge on a new word based on previous answers.
*   **Vibe:** Silence vs. Chaos.

### 2. Polarity (The Spectrum Game)
*   **Concept:** Where do we stand on subjective topics?
*   **Loop:** Statement given ("Hotdogs are sandwiches") -> Users slider 0-100% -> Guess where a specific target player landed.
*   **Vibe:** Debate & "Getting to know you".

### 3. Pecking Order (Ranking/Sorting)
*   **Concept:** How does the group perceive itself?
*   **Loop:** Prompt ("Most likely to survive zombie apocalypse") -> Everyone ranks the other players -> Reveal consensus vs. biggest disagreements.
*   **Vibe:** Roasting friends (requires established groups).

### 4. Metadata (Phone as Player)
*   **Concept:** Using the device's actual state as content.
*   **Loop:** Prompt ("Battery %", "Unread Emails") -> Users enter true stats -> Group guesses who owns the extreme values.
*   **Vibe:** Low-stakes roasting based on habits.

---

## Solo / Intelligence / Tracking (Async/Daily)
*Focus: Cognitive health, leaderboards, daily habits. Requires Auth.*

### 5. Rulebreaker (Inductive Logic)
*   **Concept:** Figure out the hidden rule by testing examples.
*   **Loop:** See sequence `2, 4, 8`. Submit test sequence `3, 6, 12`. System says "Yes/No". Guess the rule ("Double previous number").
*   **Why:** "Programmer-brain" satisfaction.

### 6. Chronology (Knowledge Sorting)
*   **Concept:** Pinpointing where things belong in time.
*   **Loop:** Draw a card ("Release of Shrek"). Place it correctly on your timeline relative to existing cards ("Before iPhone release").
*   **Data Source:** Wikipedia API for massive, accurate event dataset.
*   **Why:** Addictive, educational, high replayability.

### 7. Recall (Visual Memory)
*   **Concept:** Spatial/Visual memory test.
*   **Loop:** Grid of symbols/colors flashes for 3 seconds -> Clears -> Question: "Where was the Blue Triangle?". Levels get harder.
*   **Why:** Pure cognitive metric.

### 8. Echo (Pattern Repetition)
*   **Concept:** Increasing sequence memory (Simon Says style).
*   **Loop:** System highlights sequence (Red, Blue). Player repeats. System adds one (Red, Blue, Green). Repeat until fail.
*   **Why:** Tracks working memory capacity over time.

### 9. Interference (Inhibition Control)
*   **Concept:** Stroop Test variant.
*   **Loop:** A word appears (e.g., "RED") painted in the color Blue. Player must click the button for **Blue** (the color), not read the word.
*   **Why:** Tests reaction time and cognitive inhibition/focus.

### 10. Fluency (Verbal Agility)
*   **Concept:** Word generation under pressure.
*   **Loop:** Letter given (e.g., "F") or Category ("Animals"). Player types as many unique words as possible in 60 seconds.
*   **Why:** Classic metric for verbal processing speed and executive function.

### 11. Mind Meld Ladders (Daily Convergence)
*   **Concept:** Build a 5-7 word ladder that steadily converges on a shared theme; solo-friendly.
*   **Loop:** Daily theme -> pick words in order -> score by semantic closeness to the crowd median.
*   **Share/Archive:** Compact ladder card + emoji bar; archive shows percentile + streaks.
*   **Why:** Daily ritual with group comparison and no live multiplayer.

### 12. Thread Weave (Crowd Categories)
*   **Concept:** Cluster 8-12 words into 3-4 latent "threads" and name them.
*   **Loop:** Solve -> system compares clusters to crowd-labeled groupings.
*   **Share/Archive:** 3/4 match grid + percentile; archive tracks improvement.
*   **Why:** Connections-style consensus puzzle that works solo.

### 13. Echo Chain (Hidden Rule)
*   **Concept:** Build a word chain based on a hidden rule, then guess the rule.
*   **Loop:** Start word -> extend chain -> submit rule -> score by accuracy/time.
*   **Share/Archive:** Chain + rule guess; archive shows best rules solved.
*   **Why:** Single-player puzzle sourced from community rules.

### 14. Consensus Capsule (Daily Archive)
*   **Concept:** Fill 8 blanks with answers you think others will pick for a daily prompt.
*   **Loop:** Prompt -> submit list -> score by alignment with crowd distribution.
*   **Share/Archive:** Shareable percentile card; archive calendar with trendline.
*   **Why:** Pure solo play that still compares against the group.

### 15. Meld Rush (Arcade High Score)
*   **Concept:** Rapid-fire word associations with a score multiplier.
*   **Loop:** 60-90 second timer -> each prompt expects "closest" association -> streaks boost score.
*   **Share/Archive:** High-score board (daily/weekly/all-time) + copy-paste score line.
*   **Why:** Arcade feel, easy to replay, leaderboard friendly.
