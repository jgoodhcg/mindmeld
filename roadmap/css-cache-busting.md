# CSS Cache Busting

## Work Unit Summary

**Status:** done

**Problem/Intent:**
After deployments, users may have stale CSS cached in their browsers, causing visual bugs or broken layouts until they hard-refresh. We need a way to invalidate old CSS automatically when deploying new versions.

**Constraints:**
- Current setup: Tailwind compiles `styles/input.css` → `static/css/output.css`
- CSS is served as a static file
- Server-rendered architecture (Templ + HTMX)

**Proposed Approach:**
Add a content hash to the CSS filename at build time, then reference the hashed filename in templates.

**Open Questions:**
- Should we use a build-time hash or a version/git-commit based approach?
- How to pass the hashed filename to Templ templates?

---

## Notes

### Option A: Content Hash in Filename

Generate a file like `output.a1b2c3d4.css` where the hash is derived from the file contents.

**Pros:**
- Only invalidates when CSS actually changes
- Industry standard approach

**Cons:**
- Requires build tooling to generate the hash
- Need a way to inject the filename into templates

### Option B: Query String with Version/Commit

Keep `output.css` but reference it as `output.css?v=abc123` using git commit hash or build version.

**Pros:**
- Simpler to implement
- No filename changes needed

**Cons:**
- Some aggressive caches ignore query strings
- Less reliable than filename hashing

### Option C: Go Embed with Hash at Startup

Compute the hash when the Go server starts and expose it to templates.

**Pros:**
- No build step changes needed
- Hash computed from actual served file

**Cons:**
- Slightly more complex template integration

---

## Implementation Sketch (Option C)

1. In Go, compute MD5/SHA hash of `static/css/output.css` at startup
2. Store hash in a variable accessible to template rendering
3. In `layout.templ`, reference CSS as `/static/css/output.css?v={hash}`
4. Set appropriate cache headers (long cache with hash invalidation)

This approach requires minimal build changes and works well with the current architecture.
