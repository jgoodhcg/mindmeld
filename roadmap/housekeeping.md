# Housekeeping

## Work Unit Summary

**Status:** active

**Problem/Intent:**
Clean up build artifacts and project structure to prevent staleness issues and improve developer experience.

---

## TODO

- [ ] **Move CSS input file**: Relocate `static/css/input.css` to a more appropriate location (e.g., `styles/input.css` or root `input.css`). The `static/` directory should only contain generated/output files. Update paths in:
  - `Makefile` (css and css-watch targets)
  - `Dockerfile` (line 31)
- [ ] **Remove generated files from git**: Several generated files are tracked despite gitignore rules (added before rules existed). Run:
  ```bash
  git rm --cached static/css/output.css
  git rm --cached templates/*_templ.go
  git rm --cached internal/db/db.go internal/db/lobbies.sql.go internal/db/models.go
  ```

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
