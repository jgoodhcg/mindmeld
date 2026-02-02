# Roadmap Structure

This directory follows the AGENT_BLUEPRINT roadmap format.

## Structure

```
roadmap/
├── index.md       # Project overview and directory of work units
├── _template.md   # Starting point for new work units
├── *.md           # Individual work unit files (with frontmatter)
└── archived/      # Completed or dropped work units
```

## Work Unit Frontmatter

Every work unit file must begin with YAML frontmatter:

```yaml
---
title: "Feature Name"
status: idea | planned | active | paused | done | dropped
description: "One-line summary of what this work unit accomplishes"
tags: [area/frontend, type/feature]
priority: high | medium | low
created: YYYY-MM-DD
updated: YYYY-MM-DD
effort: XS | S | M | L | XL
depends-on: []
---
```

Required fields: `title`, `status`, `description`.

Recommended fields: `tags`, `priority`, `created`, `updated`.

## Status Definitions

| Status | Meaning |
|--------|---------|
| `idea` | Captured but not yet scoped |
| `planned` | Scoped and ready to start |
| `active` | Currently being worked on |
| `paused` | Started but blocked or deprioritized |
| `done` | Shipped and working |
| `dropped` | Decided not to pursue |

## Rules

- `roadmap/index.md` is the canonical overview and directory of work units.
- Status lives in frontmatter, not in prose.
- Update the `updated` field whenever you modify a work unit.
- Move `done` or `dropped` work units to `archived/`.
- Use consistent tag prefixes: `area/`, `type/`, `tech/`.
