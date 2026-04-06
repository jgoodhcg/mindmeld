---
title: "Cluster Delivery Log"
status: active
description: "Running implementation notes for Cluster: what worked, what missed, prompt count, and estimated token/cost usage."
created: 2026-02-10
updated: 2026-04-06
tags: [area/game, type/log]
priority: medium
---

# Cluster Delivery Log

## Purpose

Keep a running engineering record for Cluster delivery quality:
- what went right
- what went wrong and fixes
- how many prompts were needed
- rough token/cost estimates (when exact billing data is unavailable)

## Estimation Method

Exact token/cost telemetry is not exposed in this repo, so estimates are directional.

Assumptions used for rough USD estimates:
- input tokens priced at `$5 / 1M`
- output tokens priced at `$15 / 1M`
- blended midpoint used when split is unknown

`estimated_cost_usd ~= (input_tokens * 0.000005) + (output_tokens * 0.000015)`

## Running Entries

### 2026-04-06 - Cluster framing clarification + comparison-first reveal

Scope delivered:
- rewrote Cluster instructions and in-game copy so players answer for themselves rather than trying to predict the room
- removed default winner framing from reveal and repositioned centroid as the room's group center
- simplified reveal marker semantics for better non-color readability: filled player dots, outlined current-player dot, crosshair center marker, subtle labels, and explicit outlier labels
- replaced score-heavy reveal emphasis with center-distance and viewer-distance comparison readouts
- removed the temporary debrief/readout panels after manual review showed they added more UI than value

What went right:
- the product decision to treat Cluster as self-answer + compare made several UI decisions easier, especially wording and reveal emphasis
- the reveal became easier to scan once winner highlighting and extra copy blocks were removed
- the new marker semantics move the game in a more accessibility-friendly direction because meaning no longer depends on color alone

What still needs follow-up:
- prompt-axis curation is still the biggest quality lever
- two-player rounds need dedicated rehearsal feedback; outlier treatment may still be too aggressive in pair play
- an optional AI round summary may still be worth testing later, but only as flavor on top of structured reveal facts with deterministic fallback

### 2026-03-30 - Cluster playtest feedback capture

What was observed:
- players understood the basic premise and immediately started thinking about the prompt in terms of where it belongs on the coordinate plane
- one core rules question surfaced quickly during reveal: "am I answering this myself, or trying to guess where someone else put it?"
- the debrief prompt/readout was readable enough to support post-reveal discussion
- the session ended early because the tester had to leave for another meeting, so pacing and late-session stamina still need a longer rehearsal

What this suggests:
- the core interaction is legible, but the objective framing still needs to be sharper at round start and reveal
- prompt quality/review is now a clear gating task; content confidence matters before adding more mechanics
- centroid-only reveal may not be enough long-term; comparative or similarity-based insights could add depth after event-readiness work is done

Follow-up notes:
- prioritize a prompt-library review pass and wording clarity ahead of new mode work
- park "make it more interesting beyond centroid" ideas in the post-MVP roadmap rather than the event-readiness slice
- run another longer Cluster rehearsal to validate pacing, discussion energy, and end-of-game flow

### 2026-03-22 - Cluster selection stability + prompt variety

Scope delivered:
- randomized unused prompt-axis selection instead of walking the pool in deterministic created-at order
- prevented transient presence refreshes during active Cluster rounds from wiping another player's in-progress coordinate selection
- added e2e coverage that asserts the selected marker survives live updates while another player submits

What went right:
- the selection-randomization change stayed small and preserved the existing no-repeat-per-lobby rule
- the live-update fix stayed targeted to presence-triggered content refreshes, rather than disabling real-time behavior entirely
- validating the selected marker in e2e closed the gap between the user-visible bug and automated coverage

Follow-up notes:
- randomized selection improves replay variety, but it raises the quality bar on the entire prompt pool because weak prompts are now more likely to surface earlier
- future Cluster UI work should keep protecting local in-progress selection from unrelated websocket churn

### 2026-03-22 - Two-player Cluster support + viewer distance readout

Scope delivered:
- lowered Cluster's minimum start gate from 3 active players to 2
- kept centroid scoring unchanged so two-player rounds still reveal the midpoint target
- added viewer-relative `Dist from you` values in reveal standings so small groups have an immediate comparison metric
- updated automated coverage from 3-player round flow to 2-player flow with distance assertions

What went right:
- the gameplay change stayed localized to the start gate, reveal scoring helper, and standings template
- no alternate demo-only mode was needed; the real game loop now works for small-group demos and real pair play
- pairwise distance adds useful post-reveal discussion fuel even when centroid scoring ties in 2-player rounds

Follow-up notes:
- two-player rounds currently tie on centroid points by design because both players are equally far from the midpoint
- the new distance column is intended as additional reveal texture, not a replacement scoring rule

### 2026-03-22 - Production Cluster library import

Scope delivered:
- validated canonical source from `content/cluster/source`
- imported `cluster-library-v1` into production
- confirmed production schema was already current through migration `007`
- deactivated older `cluster-seed-v2` prompt/axis/pair rows to avoid mixing libraries in live selection

Production content state after import:
- 103 prompts
- 17 axis sets
- 515 prompt-axis pairs
- availability by audience:
  - Mild (10): 90 pairs
  - Polite (20): 440 pairs
  - Adults (30): 515 pairs

What went right:
- importer dry-run cleanly reported production-safe create counts with no managed existing rows
- no migration work was needed because production schema was already caught up
- deactivating the older seed content kept runtime selection deterministic and avoided silently mixing old/new pools

Follow-up notes:
- roadmap/docs were behind the actual tooling state; production import confirmed the file-first content pipeline is now operational, not speculative
- review tooling still exists as a local workflow, but production content operations currently rely on CLI import rather than an in-app admin surface

### 2026-02-10 - Cluster MVP implementation + stabilization

Scope delivered:
- Cluster game MVP (schema, seed content, routes, gameplay loop, centroid scoring, standings)
- 3-player Cluster E2E test
- websocket fix for host start-state refresh after players join/leave

What went right:
- Core MVP shipped from roadmap spec with one main implementation prompt.
- Scoring and no-repeat selection logic covered by unit tests.
- E2E regression reproduced and fixed quickly (host stale waiting state).

What went wrong:
- WebSocket join/leave events initially updated only player list, not game content.
- Result: host sometimes had to refresh to enable "START CLUSTER".
- Fixed by broadcasting content refresh trigger on player join/leave events.

Prompt count (feature thread):
- `1` main implementation prompt to build Cluster MVP
- `+` follow-up prompts for validation, migration, test iteration, and UX refinements
- running total so far in thread: `~8` prompts

Estimated token usage (very rough):
- prompt/input tokens: `~30k-80k`
- output tokens: `~40k-120k`
- total tokens: `~70k-200k`

Estimated cost (very rough):
- low estimate: `~$0.35`
- high estimate: `~$2.20`

Confidence in estimate:
- low to medium (directional only; not billing-authoritative)

### 2026-02-10 - UX feedback iteration (coordinate plane + scoring clarity)

Scope delivered:
- replaced numeric coordinate inputs with click/tap coordinate plane input
- unified submit/reveal plane rendering with shared component (`PlaneFrame`)
- moved axis labels onto the plane
- updated marker semantics (gray others, cyan self, amber winner outline, target centroid)
- added explicit scoring explanation copy in reveal UI
- added roadmap note for future text+rationale input experiment

What went right:
- the requested UX direction was implementable without backend contract changes (`x/y` remained normalized).
- E2E test adapted cleanly from field-fill to click-to-plot behavior.

What went wrong:
- no major regressions found yet; pending broader multi-browser manual check.

Prompt count delta:
- `+1` major UX feedback prompt

Estimated token/cost delta (very rough):
- tokens: `~15k-45k`
- cost: `~$0.10-$0.60`
