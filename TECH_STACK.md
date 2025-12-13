# Technical Architecture Decision Record: Mindmeld (Next.js Stack)

## 1. Core Stack
- **Language:** TypeScript (Node.js runtime)
  - *Why:* Strong typing with a familiar JS ecosystem; seamless in Next.js.
- **Framework:** Next.js (App Router)
  - *Why:* Full-stack pages/routes with built-in bundling, image/font optimizations, and SSR/SSG when needed.
- **Custom Server:** Node HTTP server + `ws`
  - *Why:* Single-port HTTP + WebSocket handling for lobby/game real-time updates.
- **Styling:** Tailwind CSS (v4, CLI)
  - *Why:* Rapid UI development with utility classes; no extra build tooling beyond the Tailwind CLI.

## 2. Data & State
- **Database:** PostgreSQL (managed in production)
  - *Why:* Reliable, scalable relational store; fits multiplayer/game/session data.
- **Schema Management:** SQL migrations (`db/migrations/*.sql`) via Node runner
  - *Why:* Plain SQL, simple tracking table; easy to run locally and in CI/deploy.
- **Session/Identity:** Signed lobby session token (httpOnly cookie)
  - *Why:* Bind `player_id` + `lobby_id` for reloads/reconnects without global auth.
- **Real-time:** `ws` library on the custom server
  - *Why:* Lightweight WebSocket handling for lobby sync, timers, and game events.

## 3. Infrastructure & Deployment
- **Host:** DigitalOcean App Platform
  - *Why:* Managed runtime, automatic HTTPS, environment variables, zero-ops deploys.
- **Runtime:** Node.js 20+
  - *Why:* Supported by App Platform; aligns with Next.js requirements.
- **Build Command:** `npm run build`
  - *Why:* Standard Next.js build pipeline.
- **Run Command:** `npm run start` (or `node server.js` once custom server is in place)
  - *Why:* Serves both HTTP and WS from one process.
- **Env Vars:** `DATABASE_URL`, `LOBBY_TOKEN_SECRET`, `PORT` (provided by DO)
  - *Why:* Separate secrets/config from code; consistent between local and prod.

## 4. Authentication (Strategy)
- **Phase 1 (MVP):** Anonymous lobby sessions
  - *Implementation:* Signed httpOnly cookie with `player_id` + `lobby_id`; display name per lobby.
- **Phase 2 (Future):** Accounts via Auth.js or hosted provider
  - *Why:* Add OAuth/email auth without changing lobby token flow; store `user_id` nullable now.

## 5. Development Tooling
- **Package Manager:** npm (lockfile checked in)
  - *Why:* Default for Next.js; matches DO build environment.
- **Migrations:** `npm run db:migrate`
  - *Why:* Runs SQL migrations with tracking table; works locally and in deploy hooks.
- **Linting:** ESLint (Next.js config)
  - *Why:* Catch common React/TS issues early.
- **Dev Server:** `npm run dev`
  - *Why:* Hot reload for rapid iteration; custom server will be added later for WS parity with prod.
