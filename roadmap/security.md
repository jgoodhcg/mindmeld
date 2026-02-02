---
title: "Application Security Hardening"
status: planned
description: "Harden headers, CSP, CSRF, and WebSocket authorization."
tags: [area/backend, type/security]
priority: medium
created: 2026-01-05
updated: 2026-02-02
effort: M
depends-on: []
---

# Application Security Hardening

## Work Unit Summary

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

---

### WebSocket Security

Current implementation uses `InsecureSkipVerify: true` for development. Before production:

**1. Authorization (High Priority)**
- Verify the connecting player is a member of the lobby before accepting the WebSocket
- Check `lobby_players` table in `handleWebSocket` before calling `hub.Register`
- Reject unauthorized connections with appropriate close code

**2. Origin Validation**
- Remove `InsecureSkipVerify: true` from `websocket.AcceptOptions`
- Configure allowed origins for production domain(s)
- Prevents cross-site WebSocket hijacking attacks

**3. Rate Limiting**
- Limit WebSocket connection attempts per IP/device token
- Consider middleware or connection-time checks
- Prevents resource exhaustion attacks

**4. Connection Limits**
- Cap maximum connections per lobby (e.g., 50)
- Cap maximum total connections server-wide
- Prevents memory exhaustion

**5. TLS/WSS**
- Ensure production deployment uses HTTPS (WSS automatic when page is HTTPS)
- Prevents eavesdropping on WebSocket traffic

**Implementation location:** `internal/server/handlers_ws.go`
