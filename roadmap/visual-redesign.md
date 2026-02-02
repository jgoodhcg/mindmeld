---
title: "Visual Redesign: 1970s Cerebral Sci-Fi"
status: idea
description: "Redesign the UI to a 1970s cerebral sci-fi visual system."
tags: [area/frontend, type/design]
priority: low
created: 2026-01-16
updated: 2026-02-02
effort: L
depends-on: []
---

# Visual Redesign: 1970s Cerebral Sci-Fi

## Work Unit Summary

**Problem/Intent:**
Shift the UI from its current gaming aesthetic (purple gradients, glass-morphism) to a 1970s cerebral science fiction vibe. The interface should feel like making a careful decision in a dimly lit research institute at 2 a.m., knowing it matters.

**Constraints:**
- Must remain mobile-first and touch-friendly (44px touch targets)
- Server-rendered architecture (Templ + HTMX) unchanged
- No film grain or slow easing animations in initial pass
- Tailwind CSS as the styling system

**Proposed Approach:**
1. Establish design system (colors, typography, base component styles)
2. Apply to screens in priority order: Lobby → Answer → Scoreboard → Home → Submit/Join
3. Validate the vibe on core screens before completing the full rollout

**Open Questions:**
- None currently

---

## Design Direction

### Core Vibe
Humanist, academic, psychological sci-fi. Technology is intimate, slow, and serious. The UI feels like a private research study or ethics lab, not a control panel. No spectacle. No gaminess. Everything suggests thought, alignment, and consequence.

### Visual Language
- Dark, low-contrast backgrounds (charcoal, deep blue-green)
- Soft radial gradients, subtle where used
- Light appears only where meaning exists
- One warm accent color (muted amber) reserved for titles or irreversible actions
- Cyan/teal glow used sparingly to indicate active cognition or connection

### Layout Principles
- Symmetry preferred
- Centered primary interactions
- Generous negative space
- Panels feel architectural, not floating
- UI elements align to invisible grids reminiscent of blueprints or academic journals
- Nothing touches screen edges unless intentional

### Anti-Goals
- No neon cyberpunk
- No gamified UI tropes
- No dopamine loops
- No busy HUDs
- No mascots
- No UI humor

---

## Color Palette

| Role | Current | Proposed | Hex |
|------|---------|----------|-----|
| Background (base) | `gray-900` | Deep charcoal | `#0d0f12` |
| Background (elevated) | `gray-800/50` | Slate blue-charcoal | `#151a1f` |
| Border/subtle | `gray-700/50` | Muted blue-gray | `#1e2730` |
| Text (primary) | white | Off-white/warm gray | `#e8e6e3` |
| Text (secondary) | `gray-400` | Cool gray | `#6b7280` |
| Accent (irreversible) | purple | Muted amber | `#d4a254` |
| Accent (cognition) | — | Cyan/teal | `#4fd1c5` |
| Danger | `#ef4444` | Muted rust | `#c45c4a` |
| Success | `#10b981` | Soft teal | `#3d9a8b` |

The purple gets retired entirely. The palette shifts colder and more muted.

---

## Typography

**Fonts:**
- **Space Mono** — headings and UI elements (1970s techno-academic feel)
- **Inter** — body text (clean, geometric, highly legible)

Both are free Google Fonts.

**Principles:**
- Text density > icon density
- Large titles, restrained body text
- No decorative fonts
- All-caps used sparingly for headings

---

## Semantic Color Meanings

**Amber (irreversible/commitment):**
- Submitting your final answer
- Starting the game as host
- Any "lock in" action

**Cyan/teal (cognition/connection):**
- When a player is actively thinking (pre-submission)
- When players are aligned/melded
- Active game phase indicators
- The "connection forming" moments

**Neutral states:**
- Waiting, idle, informational — use the cool grays

---

## Component Changes

**Drop:**
- Glass-morphism (replace with solid, architectural panels with subtle borders)
- Gradient buttons (replace with flat fills, subtle hover states)
- Large rounded corners (move from `rounded-2xl` to `rounded` or `rounded-sm`)
- Emoji status indicators (replace with text states or thin-line icons)

**Add:**
- Centered, symmetrical layouts for primary interactions
- More negative space
- Blueprint/technical grid alignment feel

**Keep:**
- Mobile-first responsive approach
- Touch-target sizing (44px)
- Real-time HTMX/WebSocket architecture
