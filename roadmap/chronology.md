# Chronology (Timeline Game)

**Status:** Planned / Idea Phase

**Goal:** A daily single-player game where users place historical events on a timeline.

## Gameplay Loop
1.  **Daily Seed:** Every 24h, a new set of 5-10 events is generated (e.g., "Invention of the Lightbulb", "Fall of Rome", "Release of 'Thriller'", "Pyramids Built").
2.  **The Board:** User sees one "anchor" event placed on a timeline.
3.  **Action:** User picks an unplaced event card and drags it to the correct spot relative to the placed events (Before/After).
4.  **Feedback:**
    *   If correct: It locks into place.
    *   If wrong: Life lost? Or just told "Earlier/Later"? (Design decision needed).
    *   *Preference:* "Strikes" system. 3 strikes and game over (or just reduced score).
5.  **Completion:** Game ends when all events are placed or strikes exceeded.
6.  **Share:** "MindMeld Chronology #42: 5/5 completed with 0 strikes."

## Tech Implementation
*   **Data Source:** Wikipedia API / Wikidata is best for structured dates.
*   **Storage:** `daily_chronology_puzzles` table containing the set of events for that date.
*   **Frontend:** HTML5 Drag and Drop API (native) or a lightweight library compatible with HTMX.
*   **Backend:** Validation logic (is date X < date Y?).

## UX / UI
*   **Vertical vs Horizontal:** Vertical timeline usually works better on mobile.
*   **Visuals:** Each card needs a title, maybe a small image (if Wikipedia thumb available), and eventually the year reveals itself.

## Riskiest Assumption
*   "Is it fun?" (Yes, board games like *Timeline* prove this).
*   "Can we automate quality content?" (Dates can be tricky—e.g., "Invention of the Wheel" is vague. Need to stick to concrete events or person lifespans).
