# Apex Memory - Frontend Layout Specification

> **Purpose**: This document provides authoritative specifications for implementing the Apex Memory frontend. An AI coding agent should be able to faithfully recreate the design by following these instructions.

---

## Table of Contents

1. [Design Philosophy](#design-philosophy)
2. [Application Architecture](#application-architecture)
3. [Route Structure](#route-structure)
4. [Layout System](#layout-system)
5. [Component Specification](#component-specification)
6. [Design Tokens](#design-tokens)
7. [State Management](#state-management)
8. [Implementation Guidelines](#implementation-guidelines)

---

## Design Philosophy

### Core Principles

1. **Progressive Disclosure** - Show essential information first, reveal details on demand
2. **Contextual Actions** - Tools and actions appear where they're needed
3. **Spatial Consistency** - Similar elements occupy similar positions across views
4. **Keyboard-First** - All actions accessible via keyboard (Cmd+K, 1-4 ratings, Esc)
5. **Mobile-Optimized** - Touch-friendly targets, collapsible UI, responsive grids

### Visual Language

- **Rounded corners** - Soft, approachable aesthetic (rounded-xl to rounded-3xl)
- **Gradient accents** - Sky-to-cyan gradients for primary emphasis
- **Subtle shadows** - Depth through shadows, not borders
- **Neutral foundation** - Slate grays allow content to shine
- **Consistent iconography** - Lucide icons throughout

---

## Application Architecture

### View Hierarchy

```
App Shell
├── TopNavBar (always visible)
│   ├── Logo
│   ├── HomeButton
│   ├── NotebooksDropdown
│   ├── SearchTrigger → CommandPalette
│   ├── NotebookStats (contextual)
│   └── ReviewLauncher
│
├── Main Content Area
│   ├── HomeDashboard (view: 'home')
│   │   ├── StatsHero
│   │   ├── StatCards
│   │   ├── WeeklyActivity
│   │   └── NotebooksGrid
│   │
│   └── NotebookView (view: 'notebook')
│       ├── NotebookSidebar (collapsible)
│       │   ├── SidebarSection: Sources
│       │   └── SidebarSection: Cards
│       │
│       └── MainContent
│           ├── SourceDetail (when source selected)
│           │   ├── SourceToolbar
│           │   ├── TabBar
│           │   └── TabContent
│           │
│           └── CardsGridView (when no source selected)
│
├── CommandPalette (overlay, Cmd+K)
│
└── FocusMode (fullscreen overlay, review session)
```

### State Flow

```
┌─────────────────────────────────────────────────────────────┐
│                        Global State                          │
├─────────────────────────────────────────────────────────────┤
│  currentView: 'home' | 'notebook'                           │
│  currentNotebook: Notebook | null                            │
│  commandPaletteOpen: boolean                                 │
│  focusMode: { active: boolean, scope: FocusScope }          │
├─────────────────────────────────────────────────────────────┤
│                     Notebook State                           │
├─────────────────────────────────────────────────────────────┤
│  selectedSource: Source | null                               │
│  sidebarCollapsed: boolean                                   │
│  sourceExpanded: boolean                                     │
│  cardsViewMode: 'all' | 'due' | 'mastered'                  │
└─────────────────────────────────────────────────────────────┘
```

---

## Route Structure

```
src/routes/
├── +layout.svelte              # Minimal shell, theme provider
├── +layout.server.ts           # Auth validation
│
├── (auth)/                     # Public auth routes
│   ├── +layout.svelte          # Centered card layout
│   ├── login/+page.svelte
│   └── register/+page.svelte
│
└── (app)/                      # Protected app routes
    ├── +layout.svelte          # App shell with TopNavBar
    ├── +layout.server.ts       # Load user, notebooks
    │
    ├── +page.svelte            # Home Dashboard (redirect target)
    │
    ├── notebooks/
    │   ├── +page.svelte        # Notebooks list (optional)
    │   └── [id]/
    │       ├── +page.svelte    # Notebook view (CardsGrid)
    │       ├── +page.server.ts # Load notebook, sources, cards
    │       └── sources/
    │           └── [sourceId]/
    │               └── +page.svelte  # Source detail view
    │
    └── review/
        └── +page.svelte        # Focus mode (can also be overlay)
```

### Navigation Patterns

| Trigger | Action | URL Change |
|---------|--------|------------|
| Click Home | Navigate to dashboard | `/` |
| Select notebook from dropdown | Navigate to notebook | `/notebooks/[id]` |
| Select source in sidebar | Show source detail | `/notebooks/[id]/sources/[sourceId]` |
| Click Review button | Open focus mode overlay | No change (overlay) |
| Press Cmd+K | Open command palette | No change (overlay) |

---

## Layout System

### Global Layout (App Shell)

```svelte
<!-- (app)/+layout.svelte -->
<script>
  let { data, children } = $props();
</script>

<div class="h-screen flex flex-col bg-slate-100">
  <TopNavBar
    user={data.user}
    notebooks={data.notebooks}
    currentNotebook={data.currentNotebook}
  />

  <main class="flex-1 flex overflow-hidden">
    {@render children()}
  </main>
</div>

{#if commandPaletteOpen}
  <CommandPalette />
{/if}

{#if focusModeActive}
  <FocusMode />
{/if}
```

### TopNavBar Layout

```
┌──────────────────────────────────────────────────────────────────┐
│ [Logo] [Home] [Notebooks ▼]          [🔍 Search ⌘K] [Stats] [⚡ Review] │
└──────────────────────────────────────────────────────────────────┘
```

**Structure:**
- Height: `py-2` (8px vertical padding)
- Background: `bg-white border-b border-slate-200`
- Content: `px-4 flex items-center justify-between`
- Left group: Logo, Home, Notebooks dropdown
- Right group: Search trigger, Stats badge (contextual), Review launcher

### Home Dashboard Layout

```
┌────────────────────────────────────────────────────────────┐
│                    max-w-6xl mx-auto p-8                   │
├────────────────────────────────────────────────────────────┤
│  Welcome Header (text-3xl font-bold)                       │
├────────────────────────────────────────────────────────────┤
│  ┌──────────────────────────────────────────────────────┐  │
│  │              STATS HERO (gradient card)              │  │
│  │  Cards due: 88    │  🔥 12 streak  │  ⭕ 85% ring   │  │
│  └──────────────────────────────────────────────────────┘  │
├────────────────────────────────────────────────────────────┤
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐          │
│  │Total    │ │Tomorrow │ │This Week│ │Best     │          │
│  │Cards    │ │Due      │ │         │ │Streak   │          │
│  │891      │ │42       │ │156      │ │28 days  │          │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘          │
│              grid-cols-4 gap-4                             │
├────────────────────────────────────────────────────────────┤
│  Weekly Activity Chart (bar chart, 7 days)                 │
├────────────────────────────────────────────────────────────┤
│  Your Notebooks                                            │
│  ┌────────────────────┐ ┌────────────────────┐            │
│  │ 🧬 Biology 101     │ │ 🇪🇸 Spanish B2     │            │
│  │ 156 cards • 87%    │ │ 412 cards • 82%    │            │
│  └────────────────────┘ └────────────────────┘            │
│              grid-cols-2 gap-4                             │
└────────────────────────────────────────────────────────────┘
```

### Notebook View Layout

```
┌────────────────────────────────────────────────────────────┐
│                        TopNavBar                           │
├────────────┬───────────────────────────────────────────────┤
│            │                                               │
│  Sidebar   │              Main Content                     │
│  w-72      │              flex-1                           │
│            │                                               │
│ ┌────────┐ │  When no source selected:                     │
│ │Sources │ │  ┌─────────────────────────────────────────┐  │
│ │  ...   │ │  │           CardsGridView                 │  │
│ └────────┘ │  │  [All] [Due (5)] [Mastered]             │  │
│            │  │  ┌─────┐ ┌─────┐ ┌─────┐                │  │
│ ┌────────┐ │  │  │Card │ │Card │ │Card │                │  │
│ │Cards   │ │  │  └─────┘ └─────┘ └─────┘                │  │
│ │  ...   │ │  └─────────────────────────────────────────┘  │
│ └────────┘ │                                               │
│            │  When source selected:                        │
│  Collapsed:│  ┌─────────────────────────────────────────┐  │
│  w-12      │  │           SourceDetail                  │  │
│            │  │  [Source] [Cards] [Summary] [Chat]      │  │
│            │  │           + SourceToolbar               │  │
│            │  └─────────────────────────────────────────┘  │
│            │                                               │
└────────────┴───────────────────────────────────────────────┘
```

**Sidebar States:**
- Expanded: `w-72` - Full content with labels
- Collapsed: `w-12` - Icons only
- Hidden: When `sourceExpanded === true`

### Focus Mode Layout (Fullscreen Overlay)

```
┌────────────────────────────────────────────────────────────┐
│ [✕]                    🧬 Biology 101           3/23 ═══   │
│                                                            │
│                                                            │
│              ┌──────────────────────────────┐              │
│              │                              │              │
│              │    What is the powerhouse    │              │
│              │       of the cell?           │              │
│              │                              │              │
│              │      👁 Tap to reveal        │              │
│              │                              │              │
│              └──────────────────────────────┘              │
│                                                            │
│     ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐           │
│     │ Again  │ │  Hard  │ │  Good  │ │  Easy  │           │
│     │  <1m   │ │  <10m  │ │   1d   │ │   4d   │           │
│     └────────┘ └────────┘ └────────┘ └────────┘           │
│                                                            │
│           Space to flip • 1-4 to rate • Esc to exit       │
└────────────────────────────────────────────────────────────┘
```

**Characteristics:**
- Fixed position, full viewport
- Dark background: `bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900`
- Center-aligned card with maximum readability
- Rating buttons appear after flip

---

## Component Specification

### Global Components

#### TopNavBar

```typescript
interface TopNavBarProps {
  user: User;
  notebooks: Notebook[];
  currentNotebook?: Notebook;
  currentView: 'home' | 'notebook';
}
```

**Composition:**
```svelte
<header class="bg-white border-b border-slate-200 px-4 py-2">
  <div class="flex items-center justify-between">
    <!-- Left group -->
    <div class="flex items-center gap-1">
      <LogoButton />
      <HomeButton active={currentView === 'home'} />
      <NotebooksDropdown
        {notebooks}
        current={currentNotebook}
        isInNotebook={currentView === 'notebook'}
      />
    </div>

    <!-- Right group -->
    <div class="flex items-center gap-3">
      <SearchTrigger />
      {#if currentView === 'notebook' && currentNotebook}
        <NotebookStats notebook={currentNotebook} />
      {/if}
      <ReviewLauncher {notebooks} />
    </div>
  </div>
</header>
```

#### NotebooksDropdown

```typescript
interface NotebooksDropdownProps {
  notebooks: Notebook[];
  current?: Notebook;
  isInNotebook: boolean;
  onSelect: (notebook: Notebook) => void;
}
```

**Behavior:**
- When in notebook view: Shows current notebook emoji + name
- When on home: Shows "Notebooks" with total due badge
- Dropdown shows all notebooks with due counts
- "Create new notebook" action at bottom

**Visual States:**
- Closed: Button with chevron-down
- Open: Overlay with notebook list
- Active notebook: `bg-sky-50 text-sky-900`

#### CommandPalette

```typescript
interface CommandPaletteProps {
  notebooks: Notebook[];
  currentNotebook?: Notebook;
  onClose: () => void;
}
```

**Structure:**
```
┌──────────────────────────────────────────────────────────┐
│ 🔍 Search in Biology 101...    [📖 This notebook] [🌐 All] │
├──────────────────────────────────────────────────────────┤
│ Quick actions                                             │
│ ┌──────────────────────────────────────────────────────┐ │
│ │ [+] Create new card                            ⌘ C   │ │
│ │ [↑] Add source                                 ⌘ S   │ │
│ │ [⚡] Start review                              ⌘ R   │ │
│ │ [📖] Switch notebook                           ⌘ N   │ │
│ └──────────────────────────────────────────────────────┘ │
├──────────────────────────────────────────────────────────┤
│ ↑↓ Navigate • ↵ Select • esc Close           ⌘K         │
└──────────────────────────────────────────────────────────┘
```

**Behavior:**
- Opens on Cmd+K (global keyboard listener)
- Closes on Escape or backdrop click
- Scope toggle: "This notebook" vs "All"
- Empty state: Quick actions
- With query: Search results (cards, sources, notebooks)

#### ReviewLauncher

```typescript
interface ReviewLauncherProps {
  notebooks: Notebook[];
  onStartReview: (scope: ReviewScope) => void;
}
```

**Structure:**
- Primary button: `bg-sky-500` with lightning icon
- Dropdown shows:
  - "Review All" (gradient hero option)
  - "By notebook" section with notebooks that have due cards

#### FocusMode

```typescript
interface FocusModeProps {
  notebooks: Notebook[];
  selectedNotebook?: Notebook;
  selectedSource?: Source;
  onClose: () => void;
}
```

**States:**
1. **Question state** - Card front, "Tap to reveal"
2. **Answer state** - Card front (dimmed), separator, card back, rating buttons
3. **Complete state** - Checkmark, "All done!" message

**Keyboard Controls:**
- Space: Flip card
- 1-4: Rate card (Again, Hard, Good, Easy)
- Escape: Exit focus mode

---

### Home Dashboard Components

#### StatsHero

```typescript
interface StatsHeroProps {
  totalDue: number;
  reviewedToday: number;
  currentStreak: number;
  averageRetention: number;
}
```

**Visual:**
- Gradient: `bg-gradient-to-br from-sky-500 to-cyan-600`
- Large due count: `text-4xl font-bold`
- Streak with flame emoji
- ProgressRing for retention percentage

#### StatCard

```typescript
interface StatCardProps {
  icon: Component;
  label: string;
  value: string | number;
}
```

**Visual:**
- Background: `bg-white rounded-2xl p-5 border border-slate-200`
- Icon + label in slate-500
- Value in `text-2xl font-bold text-slate-900`

#### WeeklyActivityChart

```typescript
interface WeeklyActivityChartProps {
  data: number[]; // 7 days
  labels: string[]; // ['Mon', 'Tue', ...]
}
```

**Visual:**
- Bar chart with proportional heights
- Past days: `bg-sky-500`
- Future days: `bg-slate-200`
- Max height based on highest value

#### NotebookCard

```typescript
interface NotebookCardProps {
  notebook: Notebook;
  onClick: () => void;
}
```

**Visual:**
- `bg-white hover:bg-slate-50 rounded-2xl p-5`
- Emoji + name + card count
- Due count (if any) + retention + streak

---

### Notebook View Components

#### NotebookSidebar

```typescript
interface NotebookSidebarProps {
  notebook: Notebook;
  sources: Source[];
  cards: Card[];
  selectedSource?: Source;
  isCollapsed: boolean;
  onSelectSource: (source: Source | null) => void;
  onToggleCollapse: () => void;
}
```

**Collapsed State (`w-12`):**
- Just icon buttons for Sources, Cards
- Expand button at top

**Expanded State (`w-72`):**
- Notebook header with emoji + name
- Collapsible sections for Sources and Cards
- Source list items
- Card filter buttons

#### SidebarSection

```typescript
interface SidebarSectionProps {
  title: string;
  icon: Component;
  count?: number;
  isOpen: boolean;
  onToggle: () => void;
  actions?: Snippet;
  children: Snippet;
}
```

**Visual:**
- Chevron rotates 90° when open
- Count badge on right
- Optional action buttons (e.g., "+" to add)

#### SourceListItem

```typescript
interface SourceListItemProps {
  source: Source;
  isSelected: boolean;
  onSelect: () => void;
}
```

**Visual:**
- Icon based on source type
- Name (truncated) + card count
- Selected state: `bg-sky-100`

#### SourceDetail

```typescript
interface SourceDetailProps {
  source: Source;
  cards: Card[];
  isExpanded: boolean;
  onClose: () => void;
  onStartReview: () => void;
  onToggleExpand: () => void;
}
```

**Structure:**
1. Header: Back button, source info, Review button, Expand toggle
2. TabBar: Source, Cards, Summary, Chat
3. SourceToolbar (type-specific, below tabs)
4. Tab content area

#### SourceToolbar

```typescript
interface SourceToolbarProps {
  source: Source;
  onGenerateCards: () => void;
}
```

**Type-Specific Controls:**
- **PDF**: Zoom, page navigation, search, TOC, download
- **YouTube**: Playback controls, timestamp, speed, bookmark
- **Audio**: Play/pause, seek, volume, speed
- **URL**: Open original, refresh, reader mode
- **Notes**: Edit, download

All include "Generate Cards" button on right.

#### SelectionToolbar

```typescript
interface SelectionToolbarProps {
  type: 'text' | 'image' | 'audio';
  position: { x: number; y: number };
  onAction: (action: string) => void;
  onClose: () => void;
}
```

**Actions by Type:**
- **Text**: Create card, Cloze deletion, Highlight, Add note
- **Image**: Image occlusion, Card with image
- **Audio**: Card from segment, Timestamp note

**Visual:**
- Floating toolbar, dark background (`bg-slate-900`)
- Positioned above selection
- Close button on right

#### CardsGridView

```typescript
interface CardsGridViewProps {
  cards: Card[];
  sources: Source[];
  viewMode: 'all' | 'due' | 'mastered';
  onViewModeChange: (mode: string) => void;
}
```

**Structure:**
1. Header: Title, count, view mode toggle, add button
2. Grid: Responsive columns (1 → 2 → 3)
3. Card items with source, due status, front/back preview
4. "Create card" placeholder button

---

### Shared UI Components

#### ProgressRing

```typescript
interface ProgressRingProps {
  progress: number; // 0-100
  size?: number;
  stroke?: number;
  trackColor?: string;
  progressColor?: string;
}
```

**Implementation:**
- SVG with two circles (track + progress)
- Rotated -90° for top start
- Animated stroke-dashoffset

#### SourceIcon

```typescript
interface SourceIconProps {
  type: 'pdf' | 'youtube' | 'url' | 'audio' | 'notes';
  class?: string;
}
```

**Icon Mapping:**
- pdf → FileText
- youtube → Youtube
- url → Link
- audio → Mic
- notes → StickyNote

---

## Design Tokens

### Color Palette

```css
/* Primary (Sky/Cyan) */
--color-primary-50: oklch(0.97 0.02 200);
--color-primary-100: oklch(0.93 0.04 200);
--color-primary-500: oklch(0.69 0.15 200);  /* Main accent */
--color-primary-600: oklch(0.60 0.17 200);  /* Hover state */

/* Gradients */
--gradient-hero: linear-gradient(to bottom right, var(--color-sky-500), var(--color-cyan-600));
--gradient-dark: linear-gradient(to bottom right, var(--color-slate-900), var(--color-slate-800), var(--color-slate-900));

/* FSRS Rating Colors */
--color-again: oklch(0.63 0.21 25);   /* Red */
--color-hard: oklch(0.75 0.18 60);    /* Amber */
--color-good: oklch(0.72 0.19 145);   /* Emerald */
--color-easy: oklch(0.69 0.15 200);   /* Sky */

/* Semantic */
--color-surface: white;
--color-surface-hover: var(--color-slate-50);
--color-border: var(--color-slate-200);
--color-text-primary: var(--color-slate-900);
--color-text-secondary: var(--color-slate-500);
--color-text-muted: var(--color-slate-400);
```

### Spacing Scale

```css
/* Used consistently throughout */
--space-1: 0.25rem;  /* 4px */
--space-2: 0.5rem;   /* 8px */
--space-3: 0.75rem;  /* 12px */
--space-4: 1rem;     /* 16px */
--space-5: 1.25rem;  /* 20px */
--space-6: 1.5rem;   /* 24px */
--space-8: 2rem;     /* 32px */
```

### Border Radius

```css
--radius-sm: 0.25rem;   /* 4px - Small buttons, badges */
--radius-md: 0.5rem;    /* 8px - Inputs, small cards */
--radius-lg: 0.75rem;   /* 12px - Buttons, list items */
--radius-xl: 1rem;      /* 16px - Cards, modals */
--radius-2xl: 1.5rem;   /* 24px - Large cards, hero sections */
--radius-3xl: 1.5rem;   /* 24px - Focus mode card */
```

### Typography

```css
/* Size scale */
--text-xs: 0.75rem;    /* 12px - Badges, timestamps */
--text-sm: 0.875rem;   /* 14px - Body, buttons */
--text-base: 1rem;     /* 16px - Default */
--text-lg: 1.125rem;   /* 18px - Subheadings */
--text-xl: 1.25rem;    /* 20px - Section titles */
--text-2xl: 1.5rem;    /* 24px - Card values */
--text-3xl: 1.875rem;  /* 30px - Page titles */
--text-4xl: 2.25rem;   /* 36px - Hero numbers */

/* Weights */
--font-normal: 400;
--font-medium: 500;
--font-semibold: 600;
--font-bold: 700;
```

### Shadows

```css
/* Elevation levels */
--shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.05);
--shadow-md: 0 4px 6px rgba(0, 0, 0, 0.1);
--shadow-lg: 0 10px 15px rgba(0, 0, 0, 0.1);
--shadow-xl: 0 20px 25px rgba(0, 0, 0, 0.15);
--shadow-2xl: 0 25px 50px rgba(0, 0, 0, 0.25);
```

### Transitions

```css
--transition-fast: 100ms ease;
--transition-default: 150ms ease;
--transition-slow: 300ms ease;
--transition-slower: 500ms ease;
```

---

## State Management

### Global State (Context/Store)

```typescript
// $lib/stores/app.svelte.ts
class AppState {
  // View state
  currentView = $state<'home' | 'notebook'>('home');

  // Notebook state
  currentNotebook = $state<Notebook | null>(null);
  selectedSource = $state<Source | null>(null);

  // UI state
  sidebarCollapsed = $state(false);
  sourceExpanded = $state(false);
  cardsViewMode = $state<'all' | 'due' | 'mastered'>('all');

  // Overlay state
  commandPaletteOpen = $state(false);
  focusMode = $state<{
    active: boolean;
    notebook: Notebook | null;
    source: Source | null;
  }>({ active: false, notebook: null, source: null });

  // Computed
  totalDue = $derived(() => {
    // Sum due counts from all notebooks
  });
}

export const appState = new AppState();
```

### Component State Patterns

**Local UI State:**
```svelte
<script>
  let dropdownOpen = $state(false);
  let searchQuery = $state('');
</script>
```

**Derived State:**
```svelte
<script>
  let { cards } = $props();
  let dueCards = $derived(cards.filter(c => c.due));
  let masteredCards = $derived(cards.filter(c => !c.due));
</script>
```

### Data Loading (Server)

```typescript
// +page.server.ts
export const load: PageServerLoad = async ({ params, locals }) => {
  const notebook = await api.getNotebook(params.id, locals.token);
  const sources = await api.getSources(params.id, locals.token);
  const cards = await api.getCards(params.id, locals.token);

  return { notebook, sources, cards };
};
```

---

## Implementation Guidelines

### Component Organization

```
src/lib/components/
├── layout/
│   ├── app-shell.svelte        # Main app wrapper
│   ├── top-nav-bar.svelte      # Global header
│   └── notebook-sidebar.svelte # Notebook sidebar
│
├── navigation/
│   ├── notebooks-dropdown.svelte
│   ├── review-launcher.svelte
│   └── search-trigger.svelte
│
├── overlays/
│   ├── command-palette.svelte
│   └── focus-mode.svelte
│
├── dashboard/
│   ├── stats-hero.svelte
│   ├── stat-card.svelte
│   ├── weekly-activity.svelte
│   └── notebook-card.svelte
│
├── notebook/
│   ├── sidebar-section.svelte
│   ├── source-list-item.svelte
│   ├── source-detail.svelte
│   ├── source-toolbar.svelte
│   ├── selection-toolbar.svelte
│   └── cards-grid-view.svelte
│
├── cards/
│   ├── card-item.svelte
│   ├── card-editor.svelte      # Future
│   └── rating-buttons.svelte
│
└── ui/
    ├── progress-ring.svelte
    ├── source-icon.svelte
    ├── badge.svelte
    └── ... (existing shadcn components)
```

### Component Template

```svelte
<!-- component-name.svelte -->
<script lang="ts">
  import { cn } from '$lib/utils';

  interface Props {
    // Required props
    data: DataType;
    // Optional props with defaults
    variant?: 'default' | 'compact';
    class?: string;
    // Event handlers
    onclick?: () => void;
    // Slots/snippets
    children?: import('svelte').Snippet;
  }

  let {
    data,
    variant = 'default',
    class: className,
    onclick,
    children
  }: Props = $props();

  // Local state
  let isOpen = $state(false);

  // Derived values
  let displayValue = $derived(/* computation */);

  // Effects (if needed)
  $effect(() => {
    // side effects
  });
</script>

<div
  class={cn(
    'base-classes',
    variant === 'compact' && 'compact-classes',
    className
  )}
  {onclick}
>
  {#if children}
    {@render children()}
  {/if}
</div>
```

### Accessibility Requirements

1. **Keyboard Navigation**
   - All interactive elements focusable
   - Logical tab order
   - Escape closes overlays
   - Arrow keys for list navigation

2. **ARIA Attributes**
   - `aria-expanded` for collapsible sections
   - `aria-selected` for active items
   - `aria-label` for icon-only buttons
   - `role="dialog"` for modals

3. **Focus Management**
   - Focus trap in modals
   - Return focus on close
   - Visible focus indicators

4. **Screen Reader Support**
   - Meaningful alt text
   - Hidden decorative elements
   - Status announcements

### Performance Considerations

1. **Lazy Loading**
   - Command palette content
   - Source previews (PDF, video)
   - Card content beyond viewport

2. **Memoization**
   - Expensive computations with `$derived`
   - Stable callback references

3. **Virtualization** (future)
   - Large card lists
   - Source lists

### Testing Approach

1. **Component Tests**
   - Render with various props
   - User interaction flows
   - Accessibility audits

2. **Integration Tests**
   - Navigation flows
   - Data loading states
   - Error handling

3. **Visual Regression**
   - Screenshot comparisons
   - Responsive layouts

---

## Appendix: Type Definitions

```typescript
interface User {
  id: string;
  email: string;
  name: string;
}

interface Notebook {
  id: string;
  name: string;
  emoji: string;
  color: string;
  dueCount: number;
  streak: number;
  totalCards: number;
  retention: number;
}

interface Source {
  id: string;
  name: string;
  type: 'pdf' | 'youtube' | 'url' | 'audio' | 'notes';
  cards: number;
  excerpt: string;
  pages?: number;
  duration?: string;
  addedAt: string;
}

interface Card {
  id: string;
  front: string;
  back: string;
  sourceId: string;
  due: boolean;
  interval: string;
  tags: string[];
}

interface GlobalStats {
  totalCards: number;
  cardsReviewedToday: number;
  currentStreak: number;
  longestStreak: number;
  averageRetention: number;
  reviewsThisWeek: number[];
  upcomingDue: {
    today: number;
    tomorrow: number;
    thisWeek: number;
  };
}

type ReviewScope =
  | { type: 'all' }
  | { type: 'notebook'; notebook: Notebook }
  | { type: 'source'; notebook: Notebook; source: Source };
```

---

## Implementation Checklist

- [ ] **Phase 1: Foundation**
  - [ ] Set up route structure with groups
  - [ ] Implement AppShell and TopNavBar
  - [ ] Create basic navigation flow

- [ ] **Phase 2: Home Dashboard**
  - [ ] StatsHero component
  - [ ] StatCard components
  - [ ] WeeklyActivityChart
  - [ ] NotebookCard grid

- [ ] **Phase 3: Notebook View**
  - [ ] NotebookSidebar with sections
  - [ ] SourceListItem components
  - [ ] CardsGridView

- [ ] **Phase 4: Source Detail**
  - [ ] SourceDetail with tabs
  - [ ] Type-specific toolbars
  - [ ] SelectionToolbar

- [ ] **Phase 5: Overlays**
  - [ ] CommandPalette
  - [ ] FocusMode review session

- [ ] **Phase 6: Polish**
  - [ ] Keyboard shortcuts
  - [ ] Animations and transitions
  - [ ] Accessibility audit
  - [ ] Responsive testing
