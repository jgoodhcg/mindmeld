# Housekeeping

## Work Unit Summary

**Status:** completed

**Problem/Intent:**
Clean up build artifacts and project structure to prevent staleness issues and improve developer experience.

---

## TODO

(none)

## Done

- **Fix Dockerfile build order**: Reordered to run templ generate, sqlc generate, then go mod tidy
- **Remove generated files from git**: Untracked `*_templ.go` and `output.css` files via `git rm --cached`
- **Move CSS input file**: Moved to `styles/input.css`

---

## Notes

**Generated files (for reference):**
- `make css` → `static/css/output.css`
- `make templ` → `templates/*_templ.go`
- `make sqlc` → `internal/db/*.go`
