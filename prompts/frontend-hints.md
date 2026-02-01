
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


## Animation Guidelines

### Transitions
- Use `transition-colors` for hover states
- Use `transition-all` for size/position changes
- Use `transition-transform` for rotation (chevron)

### Duration
- Keep transitions snappy: default Tailwind timing
- Debounce search: 300ms


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

## Testing Checklist

After each major frontend change, verify:

- [ ] TypeScript compiles without errors
- [ ] Component renders without console errors
- [ ] Responsive layout works (mobile/tablet/desktop)
- [ ] Dark mode compatibility (if implemented)
- [ ] Keyboard navigation functional
- [ ] Loading states display correctly
- [ ] Error states handled gracefully