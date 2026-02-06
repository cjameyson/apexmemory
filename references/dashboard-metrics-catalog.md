# Dashboard & Analytics Metrics Catalog

> Brainstormed 2026-02-06. Reference for Epic B implementation.

## Architectural Decisions

| Decision | Choice | Rationale |
|---|---|---|
| **Aggregation strategy** | Daily aggregation table (`daily_stats`) | Precompute daily summaries via cron or post-review trigger. Avoids scanning full reviews table on every dashboard load. |
| **API design** | Composable granular endpoints | Multiple small endpoints fetched in parallel. Enables selective loading, independent caching, and incremental dashboard build-out. |
| **Retrievability storage** | Store `current_retrievability` on cards table | Updated by daily cron alongside aggregation. Makes `AVG(retrievability)` trivial in SQL. Goes stale between updates but daily refresh is acceptable. |
| **Heatmap payload** | 365 daily buckets (52 weeks) | Acceptable payload size as `{date, count}[]`. |
| **User settings** | Separate `user_settings` table (1:1 with users) | Keeps `users` table lean for auth hot path. Structured columns for `timezone`, `day_start_hour`, `daily_goal`, etc. |

---

## Motivational Design Principles

The dashboard serves two masters: **informing** the user and **motivating** them. These can conflict. A metric that is technically accurate can be psychologically devastating. The design must account for real human behavior -- people miss days, have busy weeks, take vacations.

### The "Day-Off Problem"

Some metrics decay when you don't study. If a user takes a well-deserved weekend off and opens the app Monday to see their scores have dropped, that's a punishment for rest. This kills motivation and can cause app abandonment.

**Metric behavior categories:**

| Category | Behavior | Risk | Examples |
|---|---|---|---|
| **Decaying** | Gets worse every day you don't study | HIGH -- punishes absence | Avg retrievability, overdue count, due cards piling up |
| **Resetting** | Resets to zero on a miss | VERY HIGH -- devastating | Streak (without freeze) |
| **Snapshot** | Only changes when you study | SAFE -- stable between sessions | True retention, rating distribution, avg time per card |
| **Cumulative** | Only goes up, never down | SAFE -- always feels good | Lifetime reviews, cards mastered, total study time |
| **Rolling** | Smoothed over window, changes slowly | LOW RISK -- gradual | 30-day retention trend, weekly consistency |

### Design Rules

1. **Hero metrics should be SAFE.** The biggest, most prominent number should never punish absence. Use snapshot or cumulative metrics for hero placement.
2. **Decaying metrics need positive framing.** "45 cards ready for you" not "45 cards overdue." "Your memory is strong at 91%" not "Your memory dropped 3%."
3. **Streaks need grace.** Consider streak freezes (earned or purchased), "study days" (5/7 still counts), or momentum-style framing that degrades gradually rather than snapping to zero.
4. **Cumulative metrics are free motivation.** Lifetime totals, milestones, personal records -- these only go up and cost nothing to compute. Use them liberally.
5. **Show trajectory, not just position.** A user at 85% retention trending UP feels better than a user at 92% trending DOWN. Arrows, sparklines, and "vs last week" comparisons help.
6. **Separate "when you study" from "right now."** True retention ("92% when you study") is always stable. Predicted recall ("right now you'd remember 89%") decays. Lead with the former.

### Naming Philosophy

Metric names should be:
- **Aspirational** -- "Mastery" sounds like an achievement, "Retention" sounds clinical
- **Clear** -- A nursing student should understand it without a stats background
- **Shareable** -- "My Mastery Score is 94%" works in a screenshot; "My weighted average retrievability is 0.937" does not
- **Positive** -- Frame what they HAVE, not what they've lost

---

## Proposed API Endpoints

```
GET /v1/stats/today          -- Today's session metrics (reviews, time, ratings, new cards)
GET /v1/stats/streak         -- Current streak, longest streak, consistency %
GET /v1/stats/heatmap        -- Daily review counts for past 365 days
GET /v1/stats/cards          -- Card state distribution, mature/young/leech counts
GET /v1/stats/retention      -- True retention, predicted avg retrievability, desired vs actual
GET /v1/stats/forecast       -- Due card forecast for next 30 days
GET /v1/stats/trends         -- Historical daily aggregates (reviews, time, retention) with date range
GET /v1/stats/notebooks      -- Per-notebook summary (due, retention, mastery %, card states)
GET /v1/stats/distributions  -- Stability, difficulty, interval histograms (power users)
```

All endpoints accept optional `?notebook_id=` filter. Time-range endpoints accept `?from=&to=` or `?period=7d|30d|90d|all`.

---

## Data Sources

### Existing Schema (Already Available)

**cards table:**
- `state` (new/learning/review/relearning) -- card state distribution
- `stability`, `difficulty` -- FSRS memory state distributions
- `due` -- forecast computation
- `last_review`, `elapsed_days`, `scheduled_days` -- retrievability computation
- `reps`, `lapses` -- leech detection, lifetime counters
- `suspended_at`, `buried_until` -- exclusion filters
- `notebook_id` -- per-notebook aggregation

**reviews table:**
- `reviewed_at` -- heatmap, streaks, time-series aggregation
- `rating` (again/hard/good/easy) -- answer button distributions, retention calc
- `review_duration_ms` -- time studied
- `mode` (scheduled/practice) -- mode filtering
- `state_before`, `state_after` -- state transition analysis, true retention
- `stability_before`, `stability_after` -- stability growth tracking
- `difficulty_before`, `difficulty_after` -- difficulty drift
- `retrievability` -- R at review time (trend analysis)
- `interval_days` -- interval distribution

### New Schema Required

**`daily_stats` table** (precomputed daily aggregates):
```
user_id, stat_date, notebook_id (nullable for global),
total_reviews, scheduled_reviews, practice_reviews,
again_count, hard_count, good_count, easy_count,
new_cards_seen, total_duration_ms,
avg_retrievability (across all cards at end of day),
cards_new, cards_learning, cards_review, cards_relearning,
cards_mature, cards_suspended
```

**`user_settings` table:**
```
user_id (PK, FK to users),
timezone (text, default 'UTC'),
day_start_hour (smallint, default 4),  -- 4am = day boundary
daily_review_goal (int, nullable),
new_cards_per_day (int, default 20),
max_reviews_per_day (int, nullable),
card_order (text, default 'due_first'),
max_interval_days (int, default 36500),
leech_threshold (int, default 8),
created_at, updated_at
```

**`cards.current_retrievability` column:**
```
current_retrievability REAL  -- updated daily by cron
```

---

## Metrics Catalog

Each metric includes:
- **Display Name** -- catchy, user-facing name (with alternatives in parentheses)
- **Type** -- `M` = primarily Motivating/gamified, `I` = primarily Informative/diagnostic, `B` = Both
- **Decay** -- how this metric behaves when user doesn't study (Safe / Decaying / Resetting / Cumulative / Rolling)

---

### Tier 1: Core Dashboard (Ship First)

High-visibility, high-motivation metrics that form the main dashboard experience.

#### A. Today's Session

| # | Display Name | Alt Names | Type | Decay | Data Source | Visualization |
|---|---|---|---|---|---|---|
| 1 | **Today's Progress** | Daily Progress, Session Progress | M | Safe | Due count + reviewed count | Progress bar with "32 / 88 cards" |
| 2 | **Cards Reviewed** | Reviews Done, Cards Completed | M | Safe | `COUNT(reviews)` today | Big number in progress bar |
| 3 | **Focus Time** | Study Time, Time Invested | M | Safe | `SUM(review_duration_ms)` today | Formatted "23m 15s" |
| 4 | **New Cards Learned** | Cards Discovered, New Today | B | Safe | `COUNT(state_before='new')` today | Small number |
| 5 | **Session Accuracy** | Today's Accuracy, Hit Rate | B | Safe | `(non-again) / total` today | Percentage badge |
| 6 | **Rating Breakdown** | Response Spread, Confidence Split | I | Safe | `COUNT` by rating today | 4-segment mini bar |

**Notes:** All "today" metrics are snapshot -- they reset daily and only grow while studying. Zero is not punishing because it's a fresh start each day, not a loss.

#### B. Streaks & Consistency

| # | Display Name | Alt Names | Type | Decay | Data Source | Visualization |
|---|---|---|---|---|---|---|
| 7 | **Study Streak** | Learning Streak, Daily Streak | M | Resetting | Consecutive study days | Big number + flame icon |
| 8 | **Best Streak** | Longest Streak, Personal Record | M | Cumulative | Historical max | Secondary text "Best: 28" |
| 9 | **Study Map** | Activity Map, Learning Calendar, Commitment Map | M | Safe | `daily_stats.total_reviews` x 365 | Calendar heatmap |
| 10 | **Weekly Rhythm** | This Week, Week Pulse | M | Rolling | Days active this week / 7 | 7-dot row (filled = studied) |

**Streak Mitigation Strategies:**
- **Streak Freeze**: User earns or buys (with points) a freeze that protects one missed day. Duolingo proved this works.
- **Grace Period**: Don't break streak until 2 consecutive misses. One-off misses happen.
- **Momentum Model**: Instead of binary streak, use a "momentum" score (0-100) that decays gradually. Studying adds +15, missing a day subtracts -5. Takes 3+ days off to drop significantly. Display as "Momentum: 85" with a gradient bar.
- **"Study Days" Goal**: Let users set goal to 5/7 days. Streak counts weeks meeting the goal, not consecutive days. More realistic for adults with responsibilities.

**Recommendation:** Implement the **Momentum** model as the default, with traditional streak as an optional display. Momentum is more forgiving, more honest about real study patterns, and still motivating. "Momentum: 92" dropping to "Momentum: 87" after a day off is much less devastating than "Streak: 0."

#### C. Retention & Knowledge State

| # | Display Name | Alt Names | Type | Decay | Data Source | Visualization |
|---|---|---|---|---|---|---|
| 11 | **Mastery Score** | Knowledge Score, Memory Strength, Total Recall | B | Decaying (slow) | `AVG(cards.current_retrievability)` | Hero ring/gauge "94%" |
| 12 | **Recall Accuracy** | True Retention, Study Accuracy | I | Safe (snapshot) | Pass rate for mature cards at review time | Ring with comparison |
| 13 | **Accuracy vs Target** | Calibration, Algorithm Health | I | Safe (snapshot) | `desired_retention` vs true retention | Delta indicator |

**On "Mastery Score" (metric #11):**
This is the hero metric and the FSRS differentiator. It IS technically decaying, but the math protects us:
- With desired_retention = 0.9, avg retrievability hovers around 94-95%
- Missing ONE day drops it by ~0.1-0.3%. Imperceptible.
- Missing a WEEK might drop it 1-2%. Still reads as "93%" -- still impressive.
- It would take 30+ days of zero study to see a dramatic drop.

**Framing strategy:** Always show Mastery Score as a whole number percentage. The difference between 94.2% and 93.8% is invisible to the user. Use color thresholds: green >= 90%, yellow 80-89%, orange 70-79%, red < 70%. Most active users will always be green.

**Alternative hero metric: "Cards Mastered"** (cumulative count of cards with stability > 21 days). This NEVER decreases and gives a powerful sense of accumulation. Could pair with Mastery Score: hero area shows both "247 Cards Mastered" (cumulative) + "94% Mastery Score" (current state).

#### D. Card State Distribution

| # | Display Name | Alt Names | Type | Decay | Data Source | Visualization |
|---|---|---|---|---|---|---|
| 14 | **Learning Pipeline** | Card Stages, Memory Pipeline, Knowledge Stages | B | Safe | `cards` grouped by `state` | Donut chart |
| 15 | **Cards Mastered** | Mature Cards, Long-Term Memory | M | Cumulative | `cards WHERE stability > 21` | Highlighted donut segment + big number |
| 16 | **Total Cards** | Collection Size, Library | I | Cumulative | `COUNT(cards)` | Center of donut |

**Donut segment naming for states:**
- New -> "Unseen" or "Ready to Learn"
- Learning -> "Learning"
- Review -> "Growing" (young, S < 21d) + "Mastered" (mature, S >= 21d)
- Relearning -> "Refreshing" or "Rebuilding"

Positive framing: every state sounds productive. "Relearning" sounds like failure; "Refreshing" sounds like maintenance.

#### E. Workload Forecast

| # | Display Name | Alt Names | Type | Decay | Data Source | Visualization |
|---|---|---|---|---|---|---|
| 17 | **Week Ahead** | Study Forecast, Upcoming Reviews, Review Planner | B | Decaying | `cards.due` bucketed by day | Bar chart, today highlighted |
| 18 | **Cards Ready** | Waiting for You, Available Now | B | Decaying | Overdue + due today | Positive badge (not warning) |
| 19 | **Tomorrow's Cards** | Up Next | I | Decaying | Cards due tomorrow | Small number |

**Framing:** "67 cards ready for you" is inviting. "67 overdue cards" is anxiety-inducing. Same data, different emotion. If overdue count is high (e.g., after vacation), show a "Welcome Back" state with a catch-up plan instead of a wall of red numbers.

#### F. Notebook Scorecards

| # | Display Name | Alt Names | Type | Decay | Data Source | Visualization |
|---|---|---|---|---|---|---|
| 20 | **Notebook Overview** | My Notebooks, Study Deck | B | Mixed | Composite per-notebook stats | Card grid |
| 21 | **Notebook Mastery** | Deck Mastery, Subject Score | M | Decaying (slow) | `AVG(current_retrievability)` per notebook | Progress bar |
| 22 | **Needs Review** | Attention Needed, Priority Flag | I | Decaying | Low retention or high overdue | Subtle badge |

**Per-notebook card design:**
```
+------------------------------------------+
| [emoji] Anatomy 101                      |
| 23 cards ready  |  Mastery: 91%          |
| [=========>     ] 142 / 186 mastered     |
| [sparkline ~~~~]  last studied: today    |
+------------------------------------------+
```

---

### Tier 2: Detailed Analytics (Fast Follow)

#### G. Historical Trends

| # | Display Name | Alt Names | Type | Decay | Data Source | Visualization |
|---|---|---|---|---|---|---|
| 23 | **Daily Activity** | Reviews Over Time, Study History | B | Safe | `daily_stats.total_reviews` | Bar chart (30d) |
| 24 | **Time Invested** | Study Hours, Focus History | M | Safe | `daily_stats.total_duration_ms` | Area chart |
| 25 | **Accuracy Trend** | Retention Over Time, Recall Trend | B | Rolling | True retention by week/month | Line chart with confidence band |
| 26 | **Growth Rate** | New Cards Over Time, Learning Velocity | B | Safe | `daily_stats.new_cards_seen` | Stacked area |
| 27 | **Confidence Trend** | Rating Improvement, Answer Quality | I | Rolling | `% (good+easy)` over time | Line chart |
| 28 | **Mastery Growth** | Cards Mastered Over Time, Knowledge Growth | M | Cumulative | `daily_stats.cards_mature` cumulative | Area chart (always up!) |

**Metric #28 "Mastery Growth" is the perfect motivational chart.** It's a cumulative line that only goes up. Even if the user takes a week off, the line stays flat (not drops). When they return and study, it climbs again. It tells the story: "Look how much you've learned over time."

#### H. Review Breakdown

| # | Display Name | Alt Names | Type | Decay | Data Source | Visualization |
|---|---|---|---|---|---|---|
| 29 | **Response Pattern** | Answer Buttons, Rating Distribution | I | Safe | Reviews by rating + maturity | Grouped bar |
| 30 | **Lapse Rate** | Forgetting Rate, Miss Rate | I | Rolling | `again / total` over time | Line chart |
| 31 | **Study Mode Split** | Scheduled vs Practice | I | Safe | Mode breakdown | Segmented bar |
| 32 | **First-Impression Score** | New Card Pass Rate, Learning Ease | B | Safe | Non-again % for new cards | Percentage |

#### I. FSRS Distributions

| # | Display Name | Alt Names | Type | Decay | Data Source | Visualization |
|---|---|---|---|---|---|---|
| 33 | **Memory Depth** | Stability Spread, Memory Strength Map | I | Safe (snapshot) | `cards.stability` | Histogram (log) |
| 34 | **Challenge Spectrum** | Difficulty Spread, Complexity Map | I | Safe (snapshot) | `cards.difficulty` [1-10] | Histogram |
| 35 | **Interval Spread** | Scheduling Profile, Review Spacing | I | Safe (snapshot) | `cards.scheduled_days` | Histogram + cumulative |
| 36 | **Knowledge Distribution** | Recall Confidence Spread | I | Decaying (slow) | `cards.current_retrievability` | Histogram |
| 37 | **Notebook Strength** | Subject Depth, Deck Stability | B | Safe (snapshot) | `AVG(stability)` per notebook | Horizontal bar |

#### J. Per-Notebook Deep Dive

| # | Display Name | Alt Names | Type | Decay | Data Source | Visualization |
|---|---|---|---|---|---|---|
| 38 | **Notebook Pipeline** | Deck Stages | I | Safe | Card states per notebook | Mini stacked bars |
| 39 | **Notebook Comparison** | Subject Comparison, Side-by-Side | I | Safe | Retention, difficulty, stability | Grouped bar |
| 40 | **Coverage** | Completion, Progress | M | Cumulative | Seen at least once / total | Progress bar |
| 41 | **Activity Pulse** | Notebook Trend | B | Safe | 7-day reviews per notebook | Sparkline |

---

### Tier 3: Power User / Deep Analytics

#### K. FSRS Memory Insights

| # | Display Name | Alt Names | Type | Decay | Data Source | Visualization |
|---|---|---|---|---|---|---|
| 42 | **Memory Reinforcement Rate** | Stability Growth, Consolidation Rate | I | Safe | `S_after / S_before` trend | Line chart |
| 43 | **Difficulty Drift** | Card Hardening, Ease Trend | I | Safe | `D_after - D_before` over time | Line chart |
| 44 | **Forgetting Curve** | Memory Decay Curve, Your Memory Model | I | Safe | Actual R vs predicted | Overlay chart |
| 45 | **Effort vs Retention** | Workload Tradeoff, Cost Curve | I | Safe | Simulated reviews/day at DR values | Line chart |
| 46 | **Sweet Spot** | Optimal Retention, Efficiency Peak | I | Safe | Minimize cost/knowledge ratio | Marked point |
| 47 | **Fading Cards** | At-Risk, Needs Refresh | B | Decaying | `current_retrievability < 0.7` | Count + list |
| 48 | **Trouble Cards** | Leeches, Stuck Cards, Problem Cards | B | Safe | `lapses >= threshold` | Table |

**Note on "Leeches":** The Anki term "leech" is derogatory toward the card. "Trouble Cards" or "Stuck Cards" is gentler and implies the card needs help, not that it's parasitic.

#### L. Time & Efficiency

| # | Display Name | Alt Names | Type | Decay | Data Source | Visualization |
|---|---|---|---|---|---|---|
| 49 | **Pace** | Speed, Avg Card Time | I | Safe | `AVG(review_duration_ms)` | "8.2s / card" |
| 50 | **Pace by Confidence** | Time by Rating | I | Safe | Avg duration by rating | Grouped bar |
| 51 | **Speed Trend** | Getting Faster, Pace Over Time | B | Safe | Avg time/card over weeks | Line chart (down = good) |
| 52 | **Peak Hours** | Best Study Time, Hourly Rhythm | I | Safe | Reviews + success by hour | Dual-axis bar |
| 53 | **Total Time Invested** | Lifetime Study, Hours Logged | M | Cumulative | `SUM(duration_ms)` all time | "42 hours" |

#### M. State Transitions & Advanced

| # | Display Name | Alt Names | Type | Decay | Data Source | Visualization |
|---|---|---|---|---|---|---|
| 54 | **Card Flow** | State Transitions, Learning Flow | I | Safe | `state_before -> state_after` | Sankey or heatmap |
| 55 | **Daily Load** | Expected Reviews, Workload | I | Decaying | `SUM(1/interval)` | Number |
| 56 | **Estimated Time** | Daily Commitment, Time Budget | I | Decaying | Load * avg pace | Duration |
| 57 | **Model Accuracy** | Algorithm Fit, FSRS Calibration | I | Safe | BCELoss | Number |
| 58 | **Your Memory Curve** | Personal Forgetting Curve | I | Safe | param[20] visualization | Curve plot |

#### N. Achievements & Milestones

| # | Display Name | Alt Names | Type | Decay | Data Source | Visualization |
|---|---|---|---|---|---|---|
| 59 | **Lifetime Reviews** | Total Reviews, All-Time Reviews | M | Cumulative | `COUNT(reviews)` | Big number + milestone |
| 60 | **Cards Created** | Library Size, Content Built | M | Cumulative | `COUNT(facts)` | Big number |
| 61 | **Personal Bests** | Records, Hall of Fame | M | Cumulative | Max reviews/day, best streak, etc. | Achievement cards |
| 62 | **Milestones** | Badges, Achievements, Trophies | M | Cumulative | Threshold-based triggers | Badge/trophy grid |
| 63 | **Weekly Recap** | Week in Review, 7-Day Summary | B | Safe | Key metrics for past 7 days | Summary card |

**Milestone Ideas (all cumulative, never lost):**

| Milestone | Trigger | Badge Name Ideas |
|---|---|---|
| First review ever | 1 review | "First Step" |
| 100 reviews | Total reaches 100 | "Century" |
| 500 reviews | Total reaches 500 | "Dedicated" |
| 1,000 reviews | Total reaches 1,000 | "Thousand Strong" |
| 10,000 reviews | Total reaches 10,000 | "Memory Machine" |
| 100 cards mastered | 100 cards with S > 21d | "Scholar" |
| 500 cards mastered | 500 cards mastered | "Knowledge Builder" |
| 7-day streak | 7 consecutive days | "Week Warrior" |
| 30-day streak | 30 consecutive days | "Monthly Master" |
| 100-day streak | 100 consecutive days | "Centurion" |
| 1 hour total study | Lifetime time >= 1h | "Hour One" |
| 10 hours total study | Lifetime time >= 10h | "Time Invested" |
| 100 hours total study | Lifetime time >= 100h | "Commitment" |
| 95%+ accuracy (30d) | Rolling 30-day accuracy | "Precision" |
| 5 notebooks active | 5+ notebooks with cards | "Renaissance" |

---

## Filters & Facets

### Primary Filters (Available on All Views)

| Filter | Values | Default |
|---|---|---|
| **Time range** | Today, 7d, 30d, 90d, All time, Custom | 30d for trends, Today for session |
| **Notebook scope** | All (global) / Specific notebook | All |

### Secondary Filters (Context-Dependent)

| Filter | Values | Applicable Views |
|---|---|---|
| **Card state** | All / New / Learning / Review / Relearning | Distributions, card lists |
| **Review mode** | Scheduled / Practice / Both | Review history, trends |
| **Fact type** | Basic / Cloze / Image Occlusion | Card breakdowns |
| **Maturity** | Young (S < 21d) / Mature (S >= 21d) | Retention, distributions |

---

## Proposed Dashboard Layout

### Mobile-First (Single Column)

```
[1] Hero Banner
    - Today's Progress bar (32/88 cards)
    - Focus Time (23m)
    - Study Streak / Momentum (prominent)

[2] Mastery Score Ring
    - "94% Mastery Score"
    - "247 Cards Mastered" (cumulative, below ring)
    - Recall Accuracy small text

[3] Study Map (Heatmap)
    - Calendar heatmap (scrollable, 365 days)
    - Tap day for detail

[4] Learning Pipeline
    - Donut: Unseen / Learning / Growing / Mastered / Refreshing
    - "Mastered" count highlighted

[5] Week Ahead
    - 7-day forecast bar chart
    - "Cards Ready" count if applicable

[6] Notebook Grid
    - Scorecard per notebook
    - Mastery %, cards ready, sparkline
    - Tap to navigate

[7] Trends (expandable or tab)
    - Daily Activity bar chart
    - Accuracy Trend line
    - Mastery Growth cumulative area

[8] Deep Analytics (separate tab)
    - Memory Depth, Challenge Spectrum
    - Peak Hours, Pace
    - Achievements & Milestones
```

### Desktop (Two Column)

```
+--------------------------------------------------+
| Hero Banner (full width)                         |
| [Today's Progress bar] [Streak/Momentum] [Focus] |
+--------------------------------------------------+
| Mastery Score Ring     | Study Map (Heatmap)      |
| 94% + 247 Mastered    | 365-day calendar         |
+--------------------------------------------------+
| Learning Pipeline      | Week Ahead               |
| Donut + legend         | 7-day bar + cards ready  |
+--------------------------------------------------+
| Notebook Scorecards (grid, 2-3 cols)             |
+--------------------------------------------------+
| Daily Activity (bar)   | Mastery Growth (area)    |
| Accuracy Trend (line)  | Time Invested (area)     |
+--------------------------------------------------+
| [Insights Tab]                                   |
| Distributions | Efficiency | FSRS | Achievements |
+--------------------------------------------------+
```

---

## Summary of Key Display Names

Quick reference for the most prominent user-facing metric names:

| Internal Concept | Display Name | Why This Name |
|---|---|---|
| Avg retrievability | **Mastery Score** | Aspirational, clear, shareable. "My Mastery Score is 94%" |
| True retention | **Recall Accuracy** | Precise, clinical but understandable. Doesn't decay. |
| Cumulative mature cards | **Cards Mastered** | Only goes up. Tangible sense of accumulation. |
| Card state donut | **Learning Pipeline** | Implies forward motion, progress through stages. |
| Streak / momentum | **Study Streak** or **Momentum** | Streak is universal; Momentum is more forgiving. |
| Calendar heatmap | **Study Map** | Visual, implies a journey. Clean name. |
| Due forecast | **Week Ahead** | Planning-oriented, forward-looking, not punishing. |
| Overdue cards | **Cards Ready** | Reframes overdue as opportunity, not failure. |
| Stability distribution | **Memory Depth** | Evocative metaphor for how deeply encoded memories are. |
| Difficulty distribution | **Challenge Spectrum** | Positive framing of difficulty as challenge, not obstacle. |
| Leech cards | **Trouble Cards** | Card has a problem, not the user. Invites action. |
| At-risk cards | **Fading Cards** | Evocative, suggests they need attention, not that user failed. |
| Review duration | **Pace** | Athletic metaphor. Getting faster feels good. |
| Total study time | **Time Invested** | Implies value gained, not time spent/wasted. |
| Retrievability < 0.7 | **Fading Cards** | Gentle urgency without alarm. |
| Notebook avg retrievability | **Notebook Mastery** | Per-subject version of Mastery Score. |
| Cards seen once / total | **Coverage** | Clear, maps to real-world "have I covered the material?" |
| Cumulative mature over time | **Mastery Growth** | The perfect motivational chart -- only goes up. |
| Weekly summary | **Weekly Recap** | Familiar pattern from fitness/health apps. |

---

## Computation Notes

### Retrievability Formula (FSRS v6)
```
R = (1 + (19/81) * (elapsed_days / stability)) ^ (-param[20])

where:
  elapsed_days = (now - last_review) in fractional days
  stability = cards.stability
  param[20] = 0.1542 (default decay, per-notebook from fsrs_settings)
```

### Mastery Score (avg retrievability)
```sql
-- Simple: use precomputed column
SELECT AVG(current_retrievability) AS mastery_score
FROM app.cards
WHERE user_id = $1
  AND suspended_at IS NULL
  AND state != 'new';

-- Display as whole-number percentage: ROUND(mastery_score * 100)
-- Color thresholds: green >= 90, yellow 80-89, orange 70-79, red < 70
```

### Momentum Score (alternative to binary streak)
```
momentum = CLAMP(0, 100, previous_momentum + delta)

where:
  delta = +15 if studied today (capped at 100)
  delta = -5 per missed day (floor at 0)

-- Properties:
-- Starts at 0, takes ~7 consecutive days to reach 100
-- One missed day: 100 -> 95 (barely noticeable)
-- Three missed days: 100 -> 85 (nudge to come back)
-- Two weeks off: 100 -> 30 (significant but recoverable in ~5 days)
-- Much more forgiving than binary streak reset
```

### Streak Calculation
Requires `user_settings.timezone` and `user_settings.day_start_hour`:
```sql
-- Day boundary for user: shift by day_start_hour in user's timezone
-- A "study day" for a user in US/Eastern with day_start_hour=4 runs
-- from 4:00 AM ET to 3:59 AM ET the next calendar day.
```

### True Retention (Recall Accuracy)
```sql
-- Mature cards only (scheduled_days >= 21)
SELECT
    COUNT(*) FILTER (WHERE rating != 'again')::float / COUNT(*) AS true_retention
FROM app.reviews
WHERE state_before = 'review'
  AND scheduled_days >= 21
  AND mode = 'scheduled'
  AND reviewed_at >= $from AND reviewed_at < $to;
```

### Daily Workload Burden
```sql
SELECT SUM(1.0 / GREATEST(scheduled_days, 1)) AS daily_burden
FROM app.cards
WHERE user_id = $1 AND state = 'review' AND suspended_at IS NULL;
```

### Leech Detection (Trouble Cards)
```sql
SELECT * FROM app.cards
WHERE user_id = $1
  AND lapses >= (SELECT leech_threshold FROM app.user_settings WHERE user_id = $1)
  AND suspended_at IS NULL;
```
