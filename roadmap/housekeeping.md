# Housekeeping

## Work Unit Summary

**Status:** active

**Problem/Intent:**
Clean up build artifacts and project structure to prevent staleness issues and improve developer experience.

---

## TODO

- [x] **URGENT: Fix Dockerfile build order**: ~~Move `sqlc generate` BEFORE `go mod tidy` so `internal/db/` exists when Go resolves imports~~ → Reordered: templ generate, sqlc generate, then go mod tidy
- [x] **URGENT: Remove generated files from git**: ~~Use `git filter-repo` or BFG~~ → Untracked remaining `*_templ.go` files; history cleanup deferred (files won't appear on fresh clone)
- [x] **Move CSS input file**: ~~Relocate `static/css/input.css` to a more appropriate location~~ → Moved to `styles/input.css`
- [x] **Remove generated files from git**: Untracked via `git rm --cached`

## Done

---

## Notes

**Dockerfile ordering issue**: `go mod tidy` runs before `sqlc generate`, so `internal/db/` doesn't exist when Go tries to resolve imports. Fix: reorder to run `sqlc generate` first, then `go mod tidy`.

**Gitignore is correct**: The rules exist, files were just added before the rules:
```
/internal/db/
templates/*_templ.go
static/css/output.css
```

**Generated files:**
- `make css` → `static/css/output.css`
- `make templ` → `templates/*_templ.go`
- `make sqlc` → `internal/db/*.go`
