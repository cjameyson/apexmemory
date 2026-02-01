# Facts & Cards Browser - AI Coding Agent Prompt

## Overview

Build a **Facts & Cards Browser** for the Apex Memory spaced repetition application. This is a notebook-scoped view that displays all facts within a notebook, with filtering, sorting, pagination, bulk selection, and expandable row details.  The notebook context sidebar (left sidebar) is already implemented and should be preserved.

**Tech Stack:** SvelteKit 2 + Svelte 5 (runes), TailwindCSS 4+, TypeScript  
**Route:** `/notebooks/[notebookId]/facts`

---

## Phase 1: Data Layer & Types

### Objective
Create TypeScript types and the `+page.server.ts` loader that fetches facts with aggregated card statistics.

### Data Structure

Facts table key fields:
- `id`, `user_id`, `notebook_id`
- `fact_type`: `'basic' | 'cloze' | 'image_occlusion'`
- `content`: JSONB with `{ version: number, fields: Field[] }`
- `source_id`: optional FK to sources
- `created_at`, `updated_at`

Cards table relevant fields:
- `fact_id`, `element_id`, `state`, `due`, `suspended_at`, `buried_until`

### Required View Model

```typescript
interface FactListItem {
  id: string;
  factType: 'basic' | 'cloze' | 'image_occlusion';
  content: FactContent;
  sourceId: string | null;
  cardCount: number;      // total cards for this fact
  dueCount: number;       // cards due now (due <= now, not suspended/buried)
  tags: string[];         // extracted from content or separate table
  createdAt: Date;
  updatedAt: Date;
}

interface FactsPageData {
  notebook: NotebookSummary;
  facts: FactListItem[];
  stats: {
    totalFacts: number;
    totalCards: number;
    totalDue: number;
    byType: { basic: number; cloze: number; image_occlusion: number };
  };
  pagination: { page: number; pageSize: number; totalPages: number };
}
```

### Server Load Function Requirements

1. Accept URL search params: `?type=`, `?q=`, `?sort=`, `?order=`, `?page=`
2. Query facts with LEFT JOIN to cards for aggregated counts
3. Compute `dueCount` as cards where `due <= now()` AND `suspended_at IS NULL` AND `buried_until IS NULL`
4. Return stats aggregated across all facts (not just current page)
5. Apply pagination (default 10 per page)

### Deliverables
- [ ] `$lib/types/facts.ts` - Type definitions
- [ ] `+page.server.ts` - Load function with filtering, sorting, pagination
- [ ] SQL queries (via sqlc or inline) for facts list with card counts

### Verification
- Load function returns correctly shaped data
- Filtering by type works
- Search query filters content
- Pagination returns correct slice
- Stats reflect full dataset, not just current page

---

## Phase 2: Core Page Layout & Header

### Objective
Build the main page structure with header section containing title, stats summary, and action buttons.

### Visual Reference

```
┌─────────────────────────────────────────────────────────────────────┐
│ Facts & Cards                                    [Review (16)] [+ Create Fact] │
│ 12 facts • 27 cards • 16 due for review                                       │
│                                                                               │
│ ┌──────┐ ┌──────┐ ┌──────┐ ┌────────────┐                                    │
│ │⊞ Basic│ │T Cloze│ │⊡ Image│ │⏱ Due Today │                                  │
│ │   6   │ │   4   │ │   2   │ │    16      │                                  │
│ └──────┘ └──────┘ └──────┘ └────────────┘                                    │
└─────────────────────────────────────────────────────────────────────┘
```

### Component Structure

```
+page.svelte
├── FactsHeader.svelte
│   ├── Title + subtitle stats
│   ├── Action buttons (Review, Create)
│   └── QuickStats (type counts + due today)
└── ... (Phase 3+)
```

### Header Specifications

**Title Section:**
- "Facts & Cards" as h2, bold, text-foreground
- Subtitle: "{totalFacts} facts • {totalCards} cards • {totalDue} due for review"

**Action Buttons:**
- Review button: bg-primary, text-primary-foreground, shows due count, links to review session
- Create Fact button: bg-foreground, text-background, opens create modal (Phase 6)

**Quick Stats Bar:**
- 4 stat cards in horizontal row with gap-6
- Each has: icon in colored square (8x8), label above, count below
- Basic: Layers icon, bg-muted
- Cloze: Type icon, bg-cloze/15
- Image: Image icon, bg-warning/15
- Due Today: Clock icon, bg-info/15, count in text-info

### Deliverables
- [ ] `FactsHeader.svelte` component
- [ ] `QuickStats.svelte` component  
- [ ] Integrate into `+page.svelte`

### Verification
- Stats display correctly from server data
- Review button shows correct due count
- Responsive layout maintained
- Color scheme matches design

---

## Phase 3: Toolbar & Filtering

### Objective
Build the filter/search toolbar with type filter tabs, search input, sort dropdown, and view toggle.

### Visual Reference

```
┌─────────────────────────────────────────────────────────────────────┐
│ 🔍 Search facts...        │ All │ ⊞ Basic │ T Cloze │ ⊡ Image │    │
│                           └─────────────────────────────────────┘   │
│                                              [Recently Updated ▼] [⊞][⊞⊞] │
└─────────────────────────────────────────────────────────────────────┘
```

### Filter State Management

Use Svelte 5 runes with URL sync:

```typescript
// Derive from URL search params, update URL on change
let typeFilter = $state<'all' | 'basic' | 'cloze' | 'image_occlusion'>('all');
let searchQuery = $state('');
let sortBy = $state<'updated' | 'created' | 'cards' | 'due'>('updated');
let sortOrder = $state<'asc' | 'desc'>('desc');
let viewMode = $state<'table' | 'grid'>('table');
```

### URL Synchronization Pattern

- Use `goto()` with `replaceState: true` for filter changes
- Debounce search input (300ms)
- Read initial state from `$page.url.searchParams` in load

### Component Structure

```
FactsToolbar.svelte
├── SearchInput (left side, flex-1, max-w-md)
├── TypeFilterTabs (pill-style toggle group)
├── SortDropdown (select element)
└── ViewModeToggle (table/grid icon buttons)
```

### Specifications

**Search Input:**
- Search icon (text-muted-foreground) positioned left
- Placeholder: "Search facts..."
- Border border-border, rounded-lg
- Focus ring ring

**Type Filter Tabs:**
- Container: bg-muted, rounded-lg, p-0.5
- Active tab: bg-background, shadow-sm, text-foreground
- Inactive: text-muted-foreground, hover:text-foreground
- Each tab shows icon + label

**Sort Dropdown:**
- Options: Recently Updated, Oldest Updated, Recently Created, Most Cards, Most Due
- Native select styled with border, rounded-lg

**View Toggle:**
- Two icon buttons in pill container
- Table icon (default), Grid icon
- Active: bg-white shadow-sm

### Deliverables
- [ ] `FactsToolbar.svelte` component
- [ ] `SearchInput.svelte` component
- [ ] `TypeFilterTabs.svelte` component
- [ ] URL sync utilities in `$lib/utils/url.ts`

### Verification
- Changing filters updates URL params
- Page reload preserves filter state
- Search is debounced
- All filter combinations work correctly

---

## Phase 4: Facts Table View

### Objective
Build the main table displaying facts with expandable rows, selection, and inline actions.

### Visual Reference

```
┌────┬────────────────────────────────────────┬──────┬───────┬─────┬─────────┐
│ ☐  │ CONTENT                                │ TYPE │ CARDS │ DUE │ ACTIONS │
├────┼────────────────────────────────────────┼──────┼───────┼─────┼─────────┤
│ ☐  │ › Difference between chromatin and...  │⊞Basic│   1   │  —  │ ✏️ 🗑️ ••• │
│    │   Chromatin is loosely coiled DNA...   │      │       │     │         │
│    │   [cell-division] [dna]                │      │       │     │         │
├────┼────────────────────────────────────────┼──────┼───────┼─────┼─────────┤
│ ☐  │ › The enzyme [...] unwinds DNA by...   │T Cloze│  2   │ (1) │ ✏️ 🗑️ ••• │
│    │   [dna] [enzymes]                      │      │       │     │         │
└────┴────────────────────────────────────────┴──────┴───────┴─────┴─────────┘
```

### Component Structure

```
FactsTable.svelte
├── TableHeader (sticky, with select-all checkbox)
└── TableBody
    └── FactTableRow.svelte (for each fact)
        ├── Checkbox cell
        ├── Content cell (expandable)
        │   ├── Chevron + primary text
        │   ├── Secondary text (for basic: answer preview)
        │   └── Tags row
        ├── Type badge
        ├── Cards count
        ├── Due count (blue circle or dash)
        └── Actions (edit, delete, more)
```

### Content Display Logic

Create a utility function to extract display text from fact content:

**Basic facts:**
- Primary: front field value
- Secondary: back field value (truncated)

**Cloze facts:**
- Primary: text with `{{c1::word}}` → `[...]`
- No secondary line

**Image occlusion:**
- Primary: title field
- Secondary: "{N} regions" count

### Type Badges

| Type | Colors | Icon |
|------|--------|------|
| Basic | bg-muted text-muted-foreground | Layers |
| Cloze | bg-cloze/15 text-cloze | Type |
| Image | bg-warning/15 text-warning | Image |

### Due Count Display

- If `dueCount > 0`: Show in bg-info/15 circle with text-info
- If `dueCount === 0`: Show em-dash "—" in text-muted-foreground

### Selection State

```typescript
let selectedIds = $state<Set<string>>(new Set());

function toggleSelect(id: string) { ... }
function toggleSelectAll() { ... }
function clearSelection() { ... }
```

### Row Expansion

- Click anywhere on content cell (except checkbox) toggles expansion
- Chevron rotates 90° when expanded
- Expansion reveals card details (Phase 5)

### Deliverables
- [ ] `FactsTable.svelte` - main table container
- [ ] `FactTableRow.svelte` - individual row with expansion
- [ ] `FactTypeBadge.svelte` - reusable type badge
- [ ] `$lib/utils/fact-display.ts` - content extraction utilities

### Verification
- All fact types display correctly
- Selection works (individual and all)
- Tags display properly
- Due counts show correct styling
- Row expansion toggles smoothly

---

## Phase 5: Expanded Row & Bulk Actions

### Objective
Build the expanded row detail view showing fact preview and cards list, plus the bulk actions bar.

### Expanded Row Content

When a row is expanded, show below the main row:

```
┌────────────────────────────────────────────────────────────────────────┐
│     ┌─────────────────────────────────────────────────────────────┐   │
│     │ Fact Preview                                    Source: Ch 3 │   │
│     │ ─────────────────────────────────────────────────────────── │   │
│     │ Front: What is the powerhouse of the cell?                  │   │
│     │ Back: The mitochondria - responsible for producing ATP...   │   │
│     └─────────────────────────────────────────────────────────────┘   │
│                                                                        │
│     Cards (3)                                                          │
│     ┌─────────────────────────────────────────────────────────────┐   │
│     │ [review] c1    Due: Jan 28  Interval: 3d  Reps: 6  Lapses: 0│   │
│     │ [review] c2    Due: Jan 30  Interval: 2d  Reps: 5  Lapses: 1│   │
│     │ [new]    c3    —                                             │   │
│     └─────────────────────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────────────────────────┘
```

### Card State Badges

| State | Colors |
|-------|--------|
| new | bg-info/15 text-info |
| learning | bg-warning/15 text-warning |
| review | bg-success/15 text-success |
| relearning | bg-error/15 text-error |

For learning/relearning, append step: "learning (1)"

### Bulk Actions Bar

Appears when `selectedIds.size > 0`:

```
┌─────────────────────────────────────────────────────────────────────┐
│ 3 selected  [Clear]              [🏷️ Add Tags] [📦 Suspend] [🗑️ Delete] │
└─────────────────────────────────────────────────────────────────────┘
```

- Background: bg-primary/5, border-b border-primary/10
- Actions on right side
- Delete button: text-destructive, hover:bg-destructive/10

### Data Fetching for Expanded Rows

Options:
1. **Eager load**: Include cards array in initial facts query (heavier)
2. **Lazy load**: Fetch cards on expansion via API call

Recommend lazy loading with component-level fetch:

```typescript
// In FactTableRow.svelte
let cards = $state<Card[] | null>(null);
let isLoadingCards = $state(false);

async function loadCards() {
  if (cards !== null) return;
  isLoadingCards = true;
  cards = await fetch(`/api/facts/${fact.id}/cards`).then(r => r.json());
  isLoadingCards = false;
}

$effect(() => {
  if (isExpanded && cards === null) loadCards();
});
```

### Deliverables
- [ ] `FactExpandedContent.svelte` - expanded row content
- [ ] `CardStateBadge.svelte` - reusable state badge
- [ ] `CardsList.svelte` - cards table within expanded view
- [ ] `BulkActionsBar.svelte` - conditional bulk actions
- [ ] `/api/facts/[factId]/cards/+server.ts` - cards fetch endpoint

### Verification
- Expansion loads cards correctly
- Loading state shown while fetching
- All card states display with correct colors
- Bulk actions bar appears/disappears correctly
- Bulk action buttons are wired (can be no-op for now)

---

## Phase 6: Pagination & Empty States

### Objective
Build pagination controls and empty state displays.

### Pagination Component

```
┌─────────────────────────────────────────────────────────────────────┐
│ Showing 1 to 10 of 12 facts                    [«] [<] [1] [2] [>] [»] │
└─────────────────────────────────────────────────────────────────────┘
```

### Specifications

**Left side:** "Showing {start} to {end} of {total} facts"

**Right side pagination controls:**
- First page button (« ChevronsLeft)
- Previous button (< ChevronLeft)  
- Page number buttons (max 5 visible, current highlighted)
- Next button (> ChevronRight)
- Last page button (» ChevronsRight)

**Styling:**
- Disabled buttons: opacity-50, cursor-not-allowed
- Current page: bg-primary text-primary-foreground
- Other pages: hover:bg-muted text-muted-foreground

### Empty States

**No facts exist (empty notebook):**
```
     ┌───────────┐
     │   ⊞      │  (Layers icon in bg-muted circle)
     └───────────┘
     No facts yet
     Create your first fact to start building your knowledge base.
     
     [+ Create your first fact]
```

**No search results:**
```
     No facts found
     No facts match "search term". Try a different search term.
```

### Deliverables
- [ ] `Pagination.svelte` - reusable pagination component
- [ ] `EmptyState.svelte` - configurable empty state
- [ ] Integrate pagination with URL params

### Verification
- Pagination navigates correctly
- URL updates with page param
- Empty states show appropriate message
- Edge cases: single page, many pages

---

## Phase 7: Grid View (Optional)

### Objective
Implement alternative grid view for facts display.

### Grid Card Design

```
┌─────────────────────────────────┐
│ ☐ ⊞ Basic                  ✏️ 🗑️│
├─────────────────────────────────┤
│ What is the powerhouse of the   │
│ cell?                           │
│                                 │
│ The mitochondria - responsible  │
│ for producing ATP...            │
├─────────────────────────────────┤
│ ⊞ 1    ⏱ 1 due           📄    │
└─────────────────────────────────┘
```

### Specifications

- Grid: 1 column mobile, 2 columns md, 3 columns lg
- Card: bg-card, rounded-xl, border-border
- Header: bg-accent with checkbox, type badge, actions
- Content: truncated to 2 lines
- Footer: bg-accent with cards count, due count, source indicator

### Deliverables
- [ ] `FactsGrid.svelte` - grid container
- [ ] `FactGridCard.svelte` - individual card

### Verification
- Grid/table toggle works
- Selection state syncs between views
- Cards display all information correctly

---

## Component Design Principles

### State Management

1. **URL as source of truth** for filters/pagination
2. **Local state** for UI-only concerns (expanded rows, hover states)
3. **Derived state** using `$derived` for computed values

### Reusability

Create these shared components in `$lib/components/ui/`:
- `Badge.svelte` - generic badge with variants
- `Checkbox.svelte` - consistent checkbox styling
- `IconButton.svelte` - icon-only button with tooltip

### Accessibility

- All interactive elements focusable
- Keyboard navigation for table rows
- ARIA labels on icon buttons
- Proper heading hierarchy

### Performance

- Virtualize table if >100 rows (future enhancement)
- Debounce search input
- Lazy load expanded row content
- Use `{#key}` blocks sparingly

---

## File Structure

```
src/routes/notebooks/[notebookId]/facts/
├── +page.svelte
├── +page.server.ts
├── FactsHeader.svelte
├── FactsToolbar.svelte
├── FactsTable.svelte
├── FactTableRow.svelte
├── FactExpandedContent.svelte
├── FactsGrid.svelte (optional)
├── FactGridCard.svelte (optional)
├── BulkActionsBar.svelte
├── Pagination.svelte
└── EmptyState.svelte

src/lib/
├── types/
│   └── facts.ts
├── utils/
│   ├── fact-display.ts
│   └── url.ts
└── components/ui/
    ├── Badge.svelte
    ├── FactTypeBadge.svelte
    └── CardStateBadge.svelte
```

---

## Testing Checklist

After each phase, verify:

- [ ] TypeScript compiles without errors
- [ ] Component renders without console errors
- [ ] Responsive layout works (mobile/tablet/desktop)
- [ ] Dark mode compatibility (if implemented)
- [ ] Keyboard navigation functional
- [ ] Loading states display correctly
- [ ] Error states handled gracefully

---

## Notes

- The notebook sidebar (sources list, notebook header) is handled by a parent layout - this page only concerns the main content area
- Tags system may be a future enhancement - for now, extract from content or show placeholder
- Create Fact modal is out of scope for this prompt - wire button to emit event or navigate
- Review button should link to `/notebooks/[notebookId]/review`