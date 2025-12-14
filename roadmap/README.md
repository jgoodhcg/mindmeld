# Mindmeld Roadmap

## Vision

Mindmeld is a platform for multiplayer party games that bring people together through shared thinking. The games explore how aligned (or misaligned) players' minds are - through trivia, word association, and other social games.

**Primary goal:** Social connection
**Secondary goal:** Stats tracking and visualization

## Platform Concepts

### Lobby

A lobby is the core concept across all games:
- Created by a host, generates a unique join code
- Shareable URL contains the join code (e.g., `/trivia/ABC123`)
- Game-specific configuration set by host
- Persists until explicitly cleaned up

### Players

- No global authentication (v1) - display name per lobby
- Join via URL, enter display name
- Identity via server-issued session token (survives reloads)

## Games

| Game | Status | Description |
|------|--------|-------------|
| [Trivia](./trivia.md) | In Progress | Players submit questions for each other |
| Word Association | Future | Wavelength-style guessing |
| Puzzles | Future | Word guessing, crosswords, grouping |
| Turn-based Strategy | Future | Simple 2D strategy |

## Platform Roadmap

### Now
- Trivia MVP (see [trivia.md](./trivia.md))

### Later
- User authentication (optional accounts)
- Game selection from homepage
- Player history/stats across games
- AI-generated content
- Pre-made question/content packs
