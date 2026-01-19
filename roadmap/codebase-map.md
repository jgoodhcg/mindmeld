# Codebase Map / Auto-Generated Documentation

## Work Unit Summary

**Status:** idea

**Problem/Intent:**
This is a "vibe coded" app where deep code reading isn't the norm. We need key touchpoints - auto-generated or easily maintained docs that give a quick mental model of the system without diving into implementation details.

**Constraints:**
- Must stay current automatically (or be trivial to regenerate)
- Should be scannable in under 2 minutes
- Focused on "what exists" not "how it works"

---

## Programmatically Derivable (Preferred)

These can be auto-generated from code, preventing staleness:

### From Go Code

| Idea | Source | Output |
|------|--------|--------|
| **Routes table** | `routes.go` | `METHOD /path → handler` |
| **Handler signatures** | `handlers_*.go` | Function names + params |
| **Event types** | `internal/events/` | Event enum + payload structs |
| **Middleware chain** | `routes.go` | Ordered middleware list |
| **Server struct fields** | `server.go` | Dependencies injected |
| **Template function calls** | `handlers_*.go` | Which handlers render which templates |
| **Environment variables** | grep `os.Getenv` | Required config |
| **Error messages** | grep `http.Error` | User-facing errors |
| **WebSocket message types** | `handlers_ws.go`, `subscriber.go` | Message shapes |
| **Import graph** | `go mod graph` or AST | Package dependencies |
| **Public API surface** | exported functions | What's callable |

### From SQL/Database

| Idea | Source | Output |
|------|--------|--------|
| **Schema ERD** | `migrations/*.sql` or `models.go` | Entity relationship diagram |
| **Table definitions** | sqlc models | Fields, types, constraints |
| **Query catalog** | `*.sql` files | Named queries + params |
| **Foreign key relationships** | migrations | Table connections |
| **Indexes** | migrations | Performance-relevant indexes |

### From Templates

| Idea | Source | Output |
|------|--------|--------|
| **Template list** | `templates/*.templ` | All UI components |
| **Template hierarchy** | grep `@TemplateName` | Parent → child calls |
| **Partial dependencies** | templ files | Which partials used where |
| **Form actions** | grep `action=` | Form → endpoint mapping |
| **HTMX attributes** | grep `hx-` | Dynamic behavior points |
| **CSS classes used** | templ files | Tailwind class inventory |

### From Project Structure

| Idea | Source | Output |
|------|--------|--------|
| **File tree** | filesystem | Directory structure |
| **Package purposes** | directory names + go files | One-liner per package |
| **Lines of code** | `cloc` or `wc` | Size per package |
| **Test coverage** | `go test -cover` | Coverage % per package |
| **TODO/FIXME list** | grep comments | Outstanding work |
| **Generated files** | `.gitignore` patterns | What's built vs committed |

### From Git

| Idea | Source | Output |
|------|--------|--------|
| **Recent changes** | `git log --oneline` | What changed lately |
| **Hot files** | `git log --stat` | Most frequently modified |
| **Contributors** | `git shortlog` | Who works on what |
| **Blame summary** | `git blame` | Code age/ownership |

### From Dependencies

| Idea | Source | Output |
|------|--------|--------|
| **Direct deps** | `go.mod` | Libraries used |
| **Dep tree** | `go mod graph` | Transitive dependencies |
| **Outdated deps** | `go list -m -u all` | Update candidates |
| **License audit** | various tools | Dependency licenses |

### From Runtime (if instrumented)

| Idea | Source | Output |
|------|--------|--------|
| **Active routes** | request logs | Actually-used endpoints |
| **Error frequency** | error logs | Common failure points |
| **Response times** | timing middleware | Slow endpoints |
| **WebSocket connections** | hub stats | Real-time usage |

---

## Semi-Stable (Manual but Rarely Changes)

These need manual creation but don't go stale quickly:

| Idea | Why Stable |
|------|------------|
| **State machine diagrams** | Phases (waiting/playing/finished) change rarely |
| **Game rules explanation** | Core mechanics don't change often |
| **Architecture overview** | High-level structure is stable |
| **Deployment topology** | Infra changes infrequently |
| **Glossary of terms** | Domain language is consistent |
| **Design system** | Colors, fonts, spacing - defined once |
| **Security model** | Auth/authz approach |
| **Data flow diagrams** | Request lifecycle |

---

## Manual and Risky (Avoid or Minimize)

These go stale quickly and create maintenance burden:

| Idea | Why Risky |
|------|-----------|
| Detailed step-by-step guides | Code changes break steps |
| Screenshots | UI changes constantly |
| Inline "why" comments | Drift from reality |
| API documentation with examples | Payloads evolve |
| Performance benchmarks | Numbers change |

---

## Recommended First Passes

### Tier 1: Generate Now (High value, easy to derive)

```bash
# Routes table
grep -E "router\.(Get|Post|Put|Delete)" internal/server/routes.go

# Environment variables
grep -r "os.Getenv" cmd/ internal/

# Event types
cat internal/events/bus.go | grep "Event"

# Template list
ls templates/*.templ

# Table list
grep "type .* struct" internal/db/models.go
```

### Tier 2: Script It (Medium effort, high value)

- Routes → Markdown table with handler names
- Schema → Mermaid ERD from models.go
- Template hierarchy → Tree from grep

### Tier 3: Tooling (More effort, ongoing value)

- `go generate` step that outputs `docs/GENERATED.md`
- CI check that docs are up-to-date
- Pre-commit hook to regenerate

---

## Output Format Options

| Format | Pros | Cons |
|--------|------|------|
| **Markdown tables** | GitHub renders, simple | Limited visuals |
| **Mermaid diagrams** | GitHub renders, visual | Syntax learning curve |
| **JSON/YAML** | Machine readable, transformable | Not human-scannable |
| **HTML** | Rich formatting | Separate hosting |
| **In CLAUDE.md** | AI agents see it | Gets long |

---

## Open Questions

- Should generated docs live in `docs/` or alongside code?
- Regenerate on commit, on release, or on-demand?
- What's the right granularity? (too detailed = noise)
- Could AI agents auto-update these as they make changes?
