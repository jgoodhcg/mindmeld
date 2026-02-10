---
title: "Game Instructions"
status: ready
description: "Game rules and role-based guidance for players"
created: 2025-02-09
updated: 2025-02-09
tags: [ux, ui, onboarding]
priority: high
---

# Game Instructions

## Intent

Players need clear, context-aware instructions to understand game rules, objectives, and their specific role (host vs. participant). This reduces confusion at game start and improves social flow.

## Specification

Add game instruction screens that display before each game begins:

- **Game Description**: Brief overview of what the game is and the goal
- **How to Play**: Step-by-step rules and mechanics
- **Role-Specific Guidance**: 
  - Host controls (start game, skip, end session)
  - Participant expectations (how to answer, timing, what to expect)
- **Example Turn**: Quick walkthrough of one round so players understand the flow

Instructions should be:
- Concise (bullet points, short paragraphs)
- Context-aware (host sees host-specific actions)
- Skippable for returning players
- Persistent (viewable via help button during game)

## Validation

- [ ] Create instruction template for Trivia game
- [ ] Host sees "Start Game" CTA, participants see "Waiting for host" status
- [ ] Help button available during game play to re-view instructions
- [ ] E2E flow: Verify instructions display correctly for both host and non-host
- [ ] Visual criteria: Match 1970s sci-fi theme from visual-redesign.md

## Scope

- Not included: In-game tutorials, interactive walkthroughs, or hints
- Not included: Gamified onboarding or achievement systems
- Scope limited to: Static instruction screens before game start

## Context

- Current games: Trivia (MVP core complete)
- Visual theme: See visual-redesign.md for 1970s cerebral sci-fi aesthetic
- Player roles: Host (game owner) vs. Participants (joiners)

## Open Questions

None
