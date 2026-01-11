# Coordinates

## Work Unit Summary

**Status:** idea

**Problem/Intent:**
Build a "social alignment" game that uses a 2D coordinate system to map the group's mental models on complex topics. Unlike linear spectrum games, this allows for nuance (e.g., something can be "Smart" but "Evil"). The screen acts as the map, players plot their position, and the fun comes from seeing the "Cluster" vs. the "Outlier."

**Constraints:**
- Mobile-first input (dragging a dot on a grid)
- Large shared screen for results
- Short rounds (30-60s) to keep energy high
- Visuals must clearly show the "Centroid" (average) and individual positions

## Game Flow

1.  **Axes Setup:**
    - The screen displays a 2D graph.
    - **X-Axis:** Variable (e.g., "Useless" ↔ "Essential")
    - **Y-Axis:** Variable (e.g., "Trashy" ↔ "Classy")
    - *Future:* Host can select or randomize axes.

2.  **The Prompt:**
    - A specific item/concept appears (e.g., "Olive Garden", "Crypto", "Cargo Shorts").

3.  **Plotting Phase:**
    - Players drag a dot on their phone screen to the coordinates they believe fit the prompt.
    - Real-time feedback: "Waiting for X players..." (No positions shown yet).

4.  **The Reveal:**
    - All player dots fade onto the main screen simultaneously.
    - Animation calculates and highlights the **Centroid** (mathematical average position of the group).

5.  **Scoring & Socializing:**
    - **Consensus Score:** 0-100 points based on proximity to the Centroid.
    - **Hot Seat:** The player furthest from the Centroid (The Outlier) is highlighted.
    - **LLM Summary (Optional):** "The group thinks this is Essential but Trashy. Justin is the outlier who thinks it's Classy."

## Data Model Extensions (Draft)

- `coordinates_rounds`: Stores axes labels and the prompt.
- `coordinates_submissions`: Stores X/Y values (0.0 to 1.0 floats) for each player.
- `coordinates_stats`: Tracks cumulative distance from centroid per player (for "Normie" vs "Contrarian" awards).

## Future Ideas
- **Trend Lines:** Show how a specific player compares to the group average across *all* questions (e.g., "You are consistently more cynical than the group").
- **Team Mode:** Red team vs Blue team consistency.
