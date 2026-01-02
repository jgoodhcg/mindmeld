# Roadmap Structure

This directory contains project management documentation for Mindmeld.

## Files

| File | Purpose |
|------|---------|
| `index.md` | Canonical state: goal, current focus, active work units, key links. Keep short; no tasks. |
| `README.md` | This file. Structure and rules for the roadmap. |
| `_template.md` | Starter template for new work units. |
| `log.md` | (Optional) Append-only notes; no retroactive edits. |
| `*.md` | Work units. Each describes a discrete piece of work. |
| `archived/` | Completed or dropped work units moved here. |

## Work Unit Format

Each work unit file starts with a summary section:

```markdown
# [Title]

## Work Unit Summary

**Status:** idea | active | paused | done | dropped

**Problem/Intent:** What problem this solves or what we're trying to achieve.

**Constraints:** Known limitations, requirements, or boundaries.

**Proposed Approach:** High-level strategy for implementation.

**Open Questions:** Unresolved decisions or unknowns.
```

The rest of the file contains narrative detail about the work - design notes, data models, context, etc. Avoid checklists.

## Rules

1. **Status lives in work units.** Each work unit declares its own status. The index only links to active units.

2. **Avoid subtasks and checklists.** Keep detail as narrative notes. Checklists fragment context and become stale.

3. **One work unit at a time.** LLMs operate on a single work unit per session. Keep units focused enough to work on in isolation.

4. **Scope splits create new units.** When a work unit grows too large or spawns distinct subprojects, create new work unit files.

5. **Never delete, only archive.** When work is done or dropped, move the file to `archived/`. Keep the history.

## Catalog of Work Units

### Active

- [trivia.md](./trivia.md) - Trivia MVP

### Ideas

(none yet)

### Archived

(none yet)
