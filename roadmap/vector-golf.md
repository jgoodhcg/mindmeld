---
title: "Vector Golf (Working Title)"
status: draft
description: "Daily semantic golf game using word embeddings and pgvector."
tags: [area/game, type/feature, tech/pgvector]
priority: low
created: 2026-01-19
updated: 2026-02-07
effort: L
depends-on: []
---

# Vector Golf (Working Title)

**Concept:** A daily semantic search game where players try to find the "Perfect Trio" of words with the minimum "Semantic Gap" score.

## Core Gameplay
1.  **The Field:** Users are presented with a grid/list of **12 words** (generated daily).
2.  **The Objective:** Find the specific set of **3 words** that has the highest possible semantic similarity score (The Apex).
3.  **The Action:** Select 3 words and submit.
4.  **The Feedback (The Gap):**
    *   The system calculates the similarity of your trio (e.g., 850).
    *   The system compares it to the pre-calculated Apex score (e.g., 1000).
    *   The result is the **Gap**: `1000 - 850 = 150`.
5.  **The Score (Golf):** Your **Total Score** is the sum of the Gaps from all your guesses.
    *   *Goal:* Keep the Total Score as low as possible.
    *   *Win Condition:* Finding the Apex (Gap = 0).

## Scoring Example
*   **Target (Apex):** 1000 (Perfect Similarity)
*   **Guess 1:** You pick `{ Apple, Car, Dog }`. Score: 200. Gap: **800**.
    *   *Current Total:* 800.
*   **Guess 2:** You pick `{ Apple, Banana, Pear }`. Score: 950. Gap: **50**.
    *   *Current Total:* 850.
*   **Guess 3:** You pick `{ Apple, Banana, Cherry }` (The Apex). Score: 1000. Gap: **0**.
    *   *Final Total:* 850.

## UI / UX
*   **Layout:**
    *   **Header:** Current Total Score (Red).
    *   **Main Area:** Grid of 12 selectable words.
    *   **Sidebar/Bottom:** History Log of guesses showing the Trio and the +Gap added.
*   **Visuals:**
    *   Simple, clean, data-heavy.
    *   No "hints" about which specific words are correct—only the aggregate score tells the story.

## Technical Architecture

### 1. Data & Vector Store
*   **Database:** PostgreSQL with `pgvector` extension.
*   **Model:** `text-embedding-3-small` (OpenAI) or local ONNX (e.g., `all-MiniLM-L6-v2`) for generating embeddings.
*   **Word Source:** A pre-seeded table of ~10k common English nouns/concepts.

### 2. Daily Generation
*   Select 1 random seed word.
*   Select 11 other words with varying degrees of similarity to the seed (some close, some far).
*   **Server-Side Calc:** Check all `C(12, 3) = 220` combinations to find the absolute max score (The Apex). This is the "Par" value.

### 3. API
*   `GET /api/golf/daily`: Returns the 12 words + the encrypted Apex Score (or session ID).
*   `POST /api/golf/check`:
    *   Input: `[word_id_1, word_id_2, word_id_3]`
    *   Logic: Calc similarity, compare to Apex.
    *   Output: `{ gap: 50, is_apex: false }`

## Shareable Artifact
```text
MindMeld Golf #42
⛳️ Score: 850
Attempts: 3
📉 800 -> 50 -> 0
[Link]
```

## Risks / Open Questions
*   **Difficulty:** Is it too hard without "One word is correct" hints?
    *   *Mitigation:* The "Gap" number acts as a "Hot/Cold" indicator. A huge drop in Gap (800 -> 50) tells you that you found the right *cluster*, even if you haven't found the perfect *trio*.
*   **Subjectivity:** Does the vector model align with human intuition? (Usually yes for concrete nouns, tricky for abstract concepts).
