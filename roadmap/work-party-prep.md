# Work Party Prep (Launch Readiness)

**Goal:** Polish the experience for a work social event. Prioritize stability and low friction over new game mechanics.

## Phase 1: Stabilize (Immediate)
**Objective:** Eliminate visual glitches ("flicker") and create a smooth, app-like feel.

- [ ] **Refactor WebSocket Refreshes**:
    - Current state: `window.location.reload()` causes full page reload white flash.
    - Desired state: Use HTMX event triggers or OOB swaps to fetch new content without refreshing.
    - **Tasks**:
        - Change `broadcastRoundAdvanced` to send an HTMX trigger event (e.g., `events.EventRoundAdvanced`).
        - Change `broadcastQuestionRevealed` to send an HTMX trigger event.
        - Update frontend `game_content` container to listen for these events and perform an `hx-get` to refresh the partial.

- [ ] **Answer Status Polish**:
    - Ensure the "Who has answered" indicators update smoothly without layout shifts.

## Phase 2: Enhance (Pre-Event)
**Objective:** Remove "writer's block" friction during the question submission phase.

- [ ] **AI Question Assist**:
    - Add a "Generate Question" button to the submit form.
    - **Backend**: Implement a simple LLM handler (using platform key for now) to generate 1 trivia question + answers.
    - **Frontend**: Button triggers `hx-post` to fetch values and populate the form inputs.
    - **Constraint**: Ensure questions are "Safe for Work" and generally accessible knowledge.

## Phase 3: Expansion (Post-Event)
**Objective:** Deeply engaging solo play with shareable scoring.

- [ ] **Chronology Game**:
    - Drag-and-drop timeline sorting.
    - Wikipedia data source.
    - Daily challenge mode.
