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
status: draft | ready | active | done | dropped
description: "One-line summary of what this work unit accomplishes"
created: YYYY-MM-DD
updated: YYYY-MM-DD
tags: []
priority: high | medium | low
---
```

Required fields: `title`, `status`, `description`, `created`, `updated`, `tags`, `priority`.

## Status Definitions

| Status | Meaning |
|--------|---------|
| `draft` | Captured but still clarifying scope/spec/validation |
| `ready` | Fully specified and executable without clarification |
| `active` | Currently being worked on |
| `done` | Shipped and working |
| `dropped` | Decided not to pursue |

## Definition Of Ready

Use `ready` only when the work unit includes:

- `Intent` with clear what + why
- `Specification` that is concrete and testable
- `Validation` with explicit checks
- `Scope` boundaries
- `Context` with key file/constraint pointers
- `Open Questions` cleared or removed

## Rules

- `roadmap/index.md` is the canonical overview and directory of work units.
- Status lives in frontmatter, not in prose.
- Update the `updated` field whenever you modify a work unit.
- Move `done` or `dropped` work units to `archived/`.
- Use consistent tag prefixes: `area/`, `type/`, `tech/`.
