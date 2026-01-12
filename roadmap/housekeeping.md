# Housekeeping

## Work Unit Summary

**Status:** active

**Problem/Intent:**
Clean up build artifacts and project structure to prevent staleness issues and improve developer experience.

---

## TODO

- [x] **Move CSS input file**: ~~Relocate `static/css/input.css` to a more appropriate location~~ → Moved to `styles/input.css`
- [x] **Remove generated files from git**: Untracked via `git rm --cached`

## Done

---

## Notes

**Dockerfile is correct**: The production build already generates all derived files fresh at build time (templ, sqlc, tailwind CSS). No committed generated files are used in production.

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
