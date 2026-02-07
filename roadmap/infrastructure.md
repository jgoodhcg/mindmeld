---
title: "Infrastructure & Scaling"
status: draft
description: "Scale WebSocket and event infrastructure beyond a single server."
tags: [area/backend, type/infra]
priority: low
created: 2026-01-30
updated: 2026-02-07
effort: L
depends-on: []
---

# Infrastructure & Scaling

## Work Unit Summary

**Problem/Intent:**
Current architecture assumes single-server deployment. WebSocket connections are managed in-memory via `internal/ws/Hub`. This works for MVP but blocks horizontal scaling.

**Constraints:**
- Must maintain real-time performance for game state updates
- Should not require major application code rewrites
- Cost-effective for low-traffic periods

---

## Current State

| Component | Storage | Scaling Limitation |
|-----------|---------|-------------------|
| HTTP handlers | Stateless | None (scales horizontally) |
| Database (Postgres) | Persistent | None (managed service scales) |
| WebSocket Hub | In-memory map | **Single instance only** |
| Event Bus | In-memory | **Single instance only** |

The Hub tracks `map[lobbyCode]map[*websocket.Conn]playerID`. If a second server instance runs, it has its own empty Hub and cannot broadcast to connections on other instances.

---

## TODO

- [ ] **Research: Postgres LISTEN/NOTIFY** - Could replace in-memory event bus. Each server subscribes to channels, broadcasts locally when notified.
- [ ] **Research: Redis Pub/Sub** - Alternative to Postgres for event distribution. May have lower latency.
- [ ] **Presence in Postgres** - Store active connections in a `presence` table with heartbeat/TTL for distributed presence queries.
- [ ] **Sticky sessions** - Simpler short-term: route all requests for a lobby to the same instance (requires load balancer config).

---

## Proposed Approach (when needed)

**Phase 1: Postgres LISTEN/NOTIFY**
1. Replace `events.InMemoryBus` with a Postgres-backed bus
2. Each server instance subscribes to relevant channels
3. Events published via `NOTIFY`, received via `LISTEN`
4. Local Hub broadcasts to local connections only

**Phase 2: Distributed Presence**
1. Add `presence` table: `(lobby_id, player_id, server_id, connected_at, last_seen)`
2. Heartbeat updates `last_seen` periodically
3. Background job cleans up stale rows (server crash recovery)
4. Live lobby count queries `presence` table instead of in-memory Hub

**Alternative: Redis**
- If Postgres NOTIFY latency is insufficient, migrate to Redis Pub/Sub
- Adds operational complexity (another service to manage)

---

## Open Questions
