# Apex Memory Dashboard Metrics Catalog (Epic B)

## Scope Anchored To Current Schema
This catalog is based on data currently available in `app.cards`, `app.reviews`, `app.facts`, and `app.notebooks` from `backend/db/schema-dump.sql`.

What is strong today:
- Card scheduling state and queue data (`state`, `due`, `reps`, `lapses`, `stability`, `difficulty`).
- Rich review event history (`rating`, `review_duration_ms`, `mode`, `state_before/after`, `interval_days`, `retrievability`, `reviewed_at`).
- Notebook-level organization and FSRS target settings (`fsrs_settings`).

Current limits to design around:
- No dedicated study-session table yet (session metrics must be derived from review events/time gaps).
- `review_duration_ms` is optional, so time-based metrics need null-safe fallbacks.
- Source/chunk/chat analytics are not in schema yet, so this catalog focuses on SRS and behavior analytics.

## Tier Definitions
- **Tier 1 (Hero / Must-Have):** Daily actionable metrics shown immediately on `/home`.
- **Tier 2 (Insight / Differentiator):** Trend and diagnostic analytics that improve decisions.
- **Tier 3 (Motivation / Gamification):** Progress loops and rewards that reinforce consistency.

## Metrics (16 Total)
| Tier | Metric ID | Catchy Candidate Name(s) | Bucket | Definition (Current Schema) | Primary UI Placement |
|---|---|---|---|---|---|
| Tier 1 | `due_now` | Mission Queue, Ready-to-Review | Informative | Count cards where `state != 'new'` and `due <= now()` and not suspended/buried. | Hero KPI card + Start Review CTA |
| Tier 1 | `daily_clear_rate` | Queue Clear %, Inbox Zero | Both | `scheduled_reviews_today / (scheduled_reviews_today + due_now)` using `reviews.mode='scheduled'` for today. | Hero progress bar |
| Tier 1 | `current_streak_days` | Fireline Streak, Never Miss Chain | Motivational | Consecutive days ending today with at least 1 scheduled review (`reviews.mode='scheduled'`). | Hero badge/chip |
| Tier 1 | `retention_14d` | Memory Shield, Recall Power | Both | Last 14 days: `(good + easy) / scheduled_reviews` from `reviews.rating` where `mode='scheduled'`. | Hero ring/gauge |
| Tier 1 | `focus_minutes_today` | Focus Fuel, Brain Time | Informative | `sum(review_duration_ms)/60000` for today (scheduled + practice), null-safe to `0`. | Hero KPI card |
| Tier 1 | `new_cards_seen_today` | New Knowledge Intake, Fresh Concepts | Both | Count reviews where `state_before='new'` and `mode='scheduled'` today. | Hero KPI card |
| Tier 2 | `overdue_pressure` | Backlog Pressure, Debt Meter | Informative | Count cards where `due < now()` and `due::date < current_date`; show count and `% of active cards`. | Trend card + warning chip |
| Tier 2 | `review_load_forecast` | Tomorrow's Loadout, Queue Forecast | Informative | Time-series from `cards.due` buckets (next 1/3/7/14 days), excluding suspended/buried. | Forecast line/area chart |
| Tier 2 | `rating_mix_trend` | Answer DNA, Response Mix | Informative | Daily stacked share of `again/hard/good/easy` from `reviews.rating` over 30/90 days. | Stacked area chart |
| Tier 2 | `mature_lapse_rate` | Lapse Radar, Forgetting Leak | Both | Last 30 days: `again` rate where `state_before IN ('review','relearning')` and `mode='scheduled'`. | Trend sparkline + threshold color |
| Tier 2 | `stability_lift` | Stability Climb, Memory Growth Curve | Both | Daily avg delta `stability_after - stability_before` for scheduled reviews. | Trend line with 7-day moving avg |
| Tier 2 | `notebook_momentum` | League Table, Subject Heat | Both | Rank notebooks by weighted score: reviews (volume), retention, and streak contribution over 7/30 days. | Notebook leaderboard |
| Tier 3 | `consistency_score` | Consistency IQ, Study Rhythm | Motivational | Composite score (0-100): active days in last 14, streak length, and low missed-day penalties. | Gamified score card |
| Tier 3 | `personal_best_streak` | Record Streak, All-Time Fire | Motivational | Max historical consecutive scheduled-study days from `reviews.reviewed_at`. | Trophy card + "chase this" hint |
| Tier 3 | `mastery_milestones` | XP Track, Memory Levels | Motivational | Lifetime milestones from cumulative scheduled reviews and cards graduating to `state='review'`. | Progress rail + badges |
| Tier 3 | `speed_precision_score` | Flow State, Focus Combo | Both | Composite: reviews/minute * quality factor (`good+easy` share), trailing 7 days. | Weekly challenge card |

## Recommended Dashboard Information Architecture

### 1) Hero Row (Immediate Action)
Place Tier 1 metrics at top with a single primary CTA.
- Left: `due_now` + `daily_clear_rate` + Start Review button.
- Middle: `current_streak_days` + `new_cards_seen_today`.
- Right: `retention_14d` ring + `focus_minutes_today`.

Design goal: user knows in <5 seconds what to do now and why it matters.

### 2) Trends Section (Learning Intelligence)
Use LayerChart for 3 core charts:
- `review_load_forecast` (next workload and risk of pileups).
- `rating_mix_trend` + `mature_lapse_rate` (quality and forgetting signals).
- `stability_lift` (FSRS progress over time).

Design goal: convert review history into coaching, not just reporting.

### 3) Notebook Performance Section
- `notebook_momentum` leaderboard with filters (7d, 30d).
- Per-notebook mini stats: due now, retention, streak contribution.

Design goal: help users decide where to focus next and create healthy competition across subjects.

### 4) Motivation Section (Gamification)
- `consistency_score`, `personal_best_streak`, `mastery_milestones`, `speed_precision_score`.
- Include milestone thresholds and "next unlock" copy.

Design goal: provide reward loops that reinforce daily behavior even when queue is light.

## Suggested Rollout Order
1. **Phase B1 (fast win):** Implement all Tier 1 hero metrics + `overdue_pressure` + `review_load_forecast`.
2. **Phase B2 (analytics depth):** Add `rating_mix_trend`, `mature_lapse_rate`, `stability_lift`, `notebook_momentum`.
3. **Phase B3 (gamification):** Add all Tier 3 composites, badges, and milestone UI.

## Data Quality Notes
- For time windows, compute with user-local day boundaries (not UTC midnight).
- Guard against null `review_duration_ms` by treating null as unknown, not zero, for precision views; use explicit fallback only for top-level totals.
- Keep denominators explicit; if denominator is zero, return null/`N/A` instead of `0%`.
