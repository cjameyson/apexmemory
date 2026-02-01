# Facts & Cards Browser - Visual Specifications Reference

## Color Palette

### Type Badges
| Type | Background | Text | Icon |
|------|------------|------|------|
| Basic | `bg-muted` | `text-muted-foreground` | Layers |
| Cloze | `bg-cloze/15` | `text-cloze` | Type |
| Image | `bg-warning/15` | `text-warning` | Image |

### Card State Badges
| State | Background | Text |
|-------|------------|------|
| new | `bg-info/15` | `text-info` |
| learning | `bg-warning/15` | `text-warning` |
| review | `bg-success/15` | `text-success` |
| relearning | `bg-error/15` | `text-error` |

### Quick Stats Icons
| Stat | Background | Icon Color |
|------|------------|------------|
| Basic | `bg-muted` | `text-muted-foreground` |
| Cloze | `bg-cloze/15` | `text-cloze` |
| Image | `bg-warning/15` | `text-warning` |
| Due Today | `bg-info/15` | `text-info` |

### Action Buttons
| Button | Background | Text | Hover |
|--------|------------|------|-------|
| Review | `bg-primary` | `text-primary-foreground` | `hover:bg-primary/90` |
| Create | `bg-foreground` | `text-background` | `hover:bg-foreground/90` |

---

## Typography

### Header
- Title: `text-xl font-bold text-foreground`
- Subtitle: `text-sm text-muted-foreground`

### Quick Stats
- Label: `text-xs text-muted-foreground`
- Count: `text-sm font-semibold text-foreground`
- Due count: `text-sm font-semibold text-primary`

### Table
- Column headers: `text-xs font-medium text-muted-foreground uppercase tracking-wider`
- Content primary: `text-sm font-medium text-foreground`
- Content secondary: `text-xs text-muted-foreground`
- Tags: `text-xs text-muted-foreground` in `bg-muted` pill

---

## Spacing & Layout

### Header Section
- Container: `px-6 py-4 border-b border-border`
- Background: `bg-gradient-to-r from-accent to-background`
- Action buttons gap: `gap-2`
- Quick stats gap: `gap-6`
- Stat icon size: `w-8 h-8`

### Toolbar
- Container: `px-6 py-3 border-b border-border`
- Search input: `max-w-md` with `pl-10 pr-4 py-2`
- Filter tabs: `px-3 py-1.5 rounded-md`
- View toggle buttons: `p-2`

### Table
- Row padding: `px-4 py-3`
- Checkbox column: `w-12`
- Type column: `w-24`
- Cards column: `w-20`
- Due column: `w-20`
- Actions column: `w-28`

### Pagination
- Container: `px-6 py-3 border-t border-border`
- Page button size: `w-8 h-8`

---

## Interactive States

### Filter Tabs
```css
/* Inactive */
text-muted-foreground hover:text-foreground

/* Active */
bg-background shadow-sm text-foreground
```

### Table Rows
```css
/* Default */
hover:bg-accent transition-colors

/* Selected */
bg-primary/5
```

### Checkboxes
- Unchecked: `text-muted-foreground` (Square icon)
- Checked: `text-primary` (CheckSquare icon)

### Expandable Rows
- Chevron default: `text-muted-foreground`
- Chevron expanded: `rotate-90` transition

---

## Component Dimensions

### Due Count Circle
- Size: `w-6 h-6`
- Style: `bg-info/15 text-info rounded-full text-xs font-medium`
- Centered with flexbox

### Tags
- Padding: `px-1.5 py-0.5`
- Border radius: `rounded`
- Font: `text-xs`

### Type Badge
- Padding: `px-2 py-0.5`
- Border radius: `rounded`
- Font: `text-xs font-medium`
- Icon size: `w-3 h-3`

### Action Buttons (in row)
- Size: `p-1.5`
- Icon size: `w-4 h-4`
- Border radius: `rounded-lg`

---

## Icons Used (from Lucide)

### Navigation & Actions
- `Search` - search input
- `Plus` - create button
- `Play` - review button
- `Edit2` - edit action
- `Trash2` - delete action
- `MoreHorizontal` - more actions
- `ChevronRight` - expand indicator
- `ChevronLeft` / `ChevronRight` - pagination
- `ChevronsLeft` / `ChevronsRight` - pagination first/last

### Type Indicators
- `Layers` - Basic type
- `Type` - Cloze type
- `Image` - Image occlusion type

### Status
- `Clock` - Due today
- `Square` - unchecked
- `CheckSquare` - checked

### Layout
- `Table` - table view
- `LayoutGrid` - grid view

---

## Animation Guidelines

### Transitions
- Use `transition-colors` for hover states
- Use `transition-all` for size/position changes
- Use `transition-transform` for rotation (chevron)

### Duration
- Keep transitions snappy: default Tailwind timing
- Debounce search: 300ms

---

## Responsive Breakpoints

### Mobile (default)
- Single column grid (if grid view)
- Condensed toolbar (stack search above filters)

### Tablet (md: 768px)
- 2 column grid
- Full toolbar

### Desktop (lg: 1024px)
- 3 column grid
- All features visible

---

## Accessibility Notes

### Focus States
- Use `focus:outline-none focus:ring-2 focus:ring-ring focus:border-transparent` for inputs
- Ensure all interactive elements are focusable

### ARIA
- `aria-label` on icon-only buttons
- `aria-expanded` on expandable rows
- `role="checkbox"` with `aria-checked` on custom checkboxes

### Keyboard
- Tab navigation through all controls
- Enter/Space to toggle checkboxes
- Arrow keys for pagination (optional enhancement)
