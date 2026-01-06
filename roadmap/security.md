# Application Security Hardening

## Work Unit Summary

**Status:** planned

**Problem/Intent:**
As the application moves towards production, we need to ensure it is resilient against common web vulnerabilities. Specifically, we want to prevent Cross-Site Scripting (XSS) and ensure user sessions are secure.

**Constraints:**
- Must not break existing functionality (Trivia MVP).
- CSP must be compatible with HTMX and Templ usage.

**Proposed Approach:**
1. **Content Security Policy (CSP):** Implement a strict CSP to restrict sources of scripts, styles, and other resources.
   - Start with `Report-Only` mode to identify violations without breaking the app.
   - Transition to enforcement mode.
2. **Security Headers:** Add standard security headers (HSTS, X-Frame-Options, X-Content-Type-Options, Referrer-Policy).
3. **CSRF Protection:** Evaluate and implement CSRF protection for form submissions.

**Open Questions:**
- Do we need a nonce-based CSP for HTMX/inline scripts, or can we stick to strict file-based policies?

---

## Notes

### CSP Implementation Details
We will likely use a middleware in `internal/server/middleware.go` to set the `Content-Security-Policy` header.

**Draft Policy:**
```http
default-src 'self';
script-src 'self' 'unsafe-inline'; 
style-src 'self' 'unsafe-inline';
img-src 'self' data:;
connect-src 'self';
```
*Note: `unsafe-inline` might be needed initially for HTMX/Alpine/Tailwind unless we implement nonces.*
