# Create Fact Modal - Implementation Prompt

## Overview

Implement a modal component for creating new facts (flashcards) in the Apex Memory application. The modal supports three fact types: Basic, Cloze, and Image Occlusion. For this task, Basic and Cloze will be fully functional; Image Occlusion will display a "Coming Soon" placeholder.

**Tech Stack**: SvelteKit 2, Svelte 5 (Runes mode), TailwindCSS 4, shadcn-svelte components, sveltekit-superforms with Zod validation.

**Key Design Principles**:
- Fast, smooth, effortless user experience
- Keyboard-first with full mouse support
- Accessible (WCAG 2.1 AA compliant)
- Use existing design tokens for all colors, spacing, typography

---

## Reference Context

Refer to the following project files for schema and data structure details:
- `003_fact_cards_reviews.sql` - Facts, cards, and reviews table definitions
- Existing type definitions in the codebase for `FactContent` JSONB structure

The content JSONB structure varies by fact type. Coordinate with backend implementation for exact field names and validation requirements.

---

## Component Architecture

```
lib/components/facts/
├── create-fact-modal.svelte          # Main modal wrapper
├── fact-type-selector.svelte         # Radio-button style type selector
├── basic-fact-editor.svelte          # Editor for basic facts
├── cloze-fact-editor.svelte          # Editor for cloze facts
├── image-occlusion-placeholder.svelte # Coming soon placeholder
├── icon-input.svelte                 # Reusable input with embedded icon
└── mini-toolbar.svelte               # Formatting toolbar (placeholder for TipTap)
```

### Reusability for Edit Modal

Several components will be reused when implementing the Edit Fact modal. Design with this in mind:

- **fact-type-selector.svelte**: In edit mode, this will be read-only/disabled (fact type cannot change after creation). Add a `disabled` or `readonly` prop.
- **basic-fact-editor.svelte**: Fully reusable. Accept initial values via props.
- **cloze-fact-editor.svelte**: Fully reusable. Accept initial values via props.
- **icon-input.svelte**: Generic component, reusable anywhere.
- **mini-toolbar.svelte**: Generic component, reusable anywhere.

Components should:
- Accept `initialData` or `value` props for pre-population
- Not assume they're always in "create" mode
- Use generic prop names (`onSave` not `onCreate`) where appropriate
- Emit events or call callbacks rather than hardcoding navigation/API calls

---

## Phase 1: Modal Shell & Type Selector - COMPLETED

### Deliverables

1. **create-fact-modal.svelte** - DONE
   - Modal overlay with backdrop blur
   - Fixed positioning: 40px from viewport top, centered horizontally
   - Max height: `calc(100vh - 80px)` with internal scroll
   - Subtle border radius (use design token for `rounded-lg` equivalent)
   - Three sections: Header (with type selector), Content (scrollable), Footer

2. **Modal Header** - DONE
   - Title: "Create Fact"
   - Close button (X icon) on right
   - Type selector integrated into header (pinned above scrollable content)

3. **fact-type-selector.svelte** - DONE
   - Three options displayed horizontally as radio-button style cards
   - Each option shows: icon, label, short description
   - Selected state: accent border, tinted background
   - Keyboard navigation (arrow keys) with `role="radiogroup"`
   - Options:
     - Basic: "Front & back"
     - Cloze: "Fill in the blank"
     - Image Occlusion: "Hide image regions"

4. **Modal Footer** - DONE
   - Left: Keyboard shortcut hints (`⌘S` save, `⌘⇧S` save & new, `Esc` cancel)
   - Right: Card count indicator, "Save & New" button, "Create Fact" primary button

5. **Props Interface** - DONE
   ```typescript
   interface CreateFactModalProps {
     open: boolean;
     notebookId: string;
     onclose: () => void;
     onsubmit: (data: FactFormData) => Promise<void>;
   }
   ```

### Exit Criteria
- [x] Modal opens/closes with proper backdrop
- [x] Modal is positioned 40px from top, doesn't shift when content changes
- [x] Type selector allows switching between all three types
- [x] Selected type is visually distinct
- [x] Footer displays static shortcut hints and buttons
- [x] Close button and Cancel button close the modal
- [x] Escape key closes the modal
- [x] Focus is trapped within modal when open
- [x] `aria-modal="true"` and proper `role="dialog"` attributes

---

## Phase 2: Basic Fact Editor - COMPLETED

### Deliverables

1. **basic-fact-editor.svelte** - DONE
   - Four fields in order: Front, Back, Back Extra, Hint
   - Props: `initialData` (optional, for edit mode), `onchange` callback, `errors`
   - Uses `cn()` (tailwind-merge) for conditional error styling
   - Design for reuse in edit modal (will receive existing fact data)

2. **Front Field** - DONE
   - Label: "Front"
   - mini-toolbar aligned right of label (placeholder buttons for now)
   - Textarea with border
   - Placeholder: "Question or prompt..."
   - Auto-focus on modal open when Basic type selected

3. **Back Field** - DONE
   - Label: "Back"
   - mini-toolbar aligned right of label
   - Textarea with border
   - Placeholder: "Answer..."

4. **Back Extra Field** (using icon-input) - DONE
   - Embedded Info icon on left inside input
   - Resizable textarea with `min-h-9` to prevent icon bleed
   - Placeholder: "Additional info shown on the back..."
   - No separate label (icon serves as visual indicator)

5. **Hint Field** (using icon-input) - DONE
   - Embedded Lightbulb icon on left inside input
   - Single-line input
   - Placeholder: "Hint (shown on request)..."
   - No separate label

6. **icon-input.svelte** (reusable) - DONE
   - Props: `icon`, `placeholder`, `value`, `multiline`, `resizable`, `oninput`, `error`
   - Icon positioned absolutely inside input padding
   - Uses `cn()` for conditional error styling (red border + red focus ring)
   - Resizable variant has `min-h-9` to prevent icon overflow

7. **mini-toolbar.svelte** (placeholder) - DONE
   - Row of icon buttons: Bold, Italic, Underline | Code, Link | Bullet List, Numbered List
   - Buttons are non-functional placeholders (will integrate TipTap later)
   - Appropriate spacing and dividers between groups

8. **State Management** - DONE
   - Each fact type maintains independent state
   - Switching types preserves data in each type's state
   - Only active type's data is submitted on save

### Exit Criteria
- [x] All four fields render with correct layout
- [x] Front field auto-focuses when Basic is selected
- [x] icon-input component works for both single-line and multiline
- [x] mini-toolbar renders (buttons can be non-functional)
- [x] Typing in fields updates component state
- [x] Switching to Cloze and back preserves Basic form data
- [x] Card count shows "1 card" for Basic type
- [x] Tab order flows logically through fields

### Human Review Checkpoint
Verify field layout, focus behavior, and state persistence before proceeding.

---

## Phase 3: Cloze Fact Editor - COMPLETED

### Deliverables

1. **cloze-fact-editor.svelte** - DONE
   - Two fields: Cloze Text, Back Extra
   - Props: `initialData` (optional, for edit mode), `onchange` callback, `oncardcount`, `errors`
   - Uses `cn()` for conditional error styling
   - Design for reuse in edit modal (will receive existing fact data)

2. **Cloze Text Field** - DONE
   - Label: "Cloze Text"
   - Toolbar row containing:
     - mini-toolbar (left)
     - Divider
     - Cloze buttons: `c1`, `c2`, `c3` (dynamically expands) + `+` button for next available
   - Monospace font textarea
   - Placeholder: "The {{c1::answer}} is hidden..."
   - Ctrl+1-9 keyboard shortcuts for cloze insertion

3. **Cloze Insertion Logic** - DONE
   - Track which cloze numbers are used in current text via `$derived`
   - Used numbers: secondary variant button with muted text
   - Unused numbers: outline variant button with primary text
   - `+` button: inserts next unused cloze number
   - Insertion behavior:
     - If text selected: wrap selection in `{{cN::selection}}`
     - If no selection: insert `{{cN::}}` and place cursor between `::` and `}}`
   - After insertion, return focus to textarea at appropriate cursor position

4. **Back Extra Field** - DONE
   - Same icon-input as Basic editor (resizable variant)
   - Placeholder: "Additional info shown on the back..."

5. **Card Count Calculation** - DONE
   - Parse cloze text for unique cloze numbers via `$derived`
   - Display count: "{N} cards" in footer
   - Regex: `/\{\{c(\d+)::/g` to find cloze markers

### Exit Criteria
- [x] Cloze text field has monospace font
- [x] Cloze buttons show correct used/unused states
- [x] Clicking cloze button wraps selected text correctly
- [x] Clicking cloze button with no selection inserts template
- [x] Cursor position is correct after insertion
- [x] `+` button inserts next available number
- [x] Card count updates dynamically as cloze markers are added/removed
- [x] Back Extra field works identically to Basic editor
- [x] Tab order is logical

---

## Phase 4: Image Occlusion Placeholder - COMPLETED

### Deliverables

1. **image-occlusion-placeholder.svelte** - DONE
   - Centered content area with:
     - Large Image icon (muted)
     - Heading: "Image Occlusion"
     - Subheading: "Coming Soon"
     - Brief description about uploading images and drawing regions
   - Muted/disabled visual treatment
   - Fills the content area of the modal

2. **Disabled State** - DONE
   - When Image Occlusion is selected, "Create Fact" and "Save & New" are disabled
   - Card count shows 0 (hidden when 0)

### Exit Criteria
- [x] Placeholder displays when Image Occlusion type selected
- [x] Create Fact button is disabled
- [x] Visual design is consistent with overall modal
- [x] User understands this is a future feature

---

## Phase 5: Form Validation & Submission - COMPLETED

Note: Validation is implemented inline in `create-fact-modal.svelte` rather than via Zod schemas. This keeps things simple for now; Zod can be introduced later if needed.

### Deliverables

1. **Inline Validation** - DONE (implemented without Zod)
   - Basic: front and back required
   - Cloze: text required, must contain at least one valid cloze deletion with content (`/\{\{c\d+::.+?\}\}/`)
   - Empty cloze `{{c1::}}` correctly rejected

2. **Validation Behavior** - DONE
   - Validate on submit attempt
   - Show inline error messages below invalid fields (red text)
   - Error border + error focus ring via `cn()` / tailwind-merge
   - Focus first invalid field on validation failure
   - Clear field errors when user starts typing (errors reset on data change)

3. **Submit Flow** - DONE
   - "Create Fact" button: validate -> submit -> close modal
   - "Save & New" button: validate -> submit -> reset form (stay on same type)
   - Loading state on buttons via `submitting` state (buttons disabled)
   - Error handling: displays inline error message above footer

4. **Form Reset** - DONE
   - After "Save & New", clear only the active type's fields
   - Preserve the selected fact type
   - `formKey` increment destroys/recreates editor components to reset internal state
   - Return focus to first field via `requestAnimationFrame`

### Exit Criteria
- [x] Empty required fields show validation errors
- [x] Cloze text without any cloze markers shows error
- [x] Empty cloze deletions `{{c1::}}` show error
- [x] Errors clear when user edits field
- [x] Submit shows loading state
- [x] "Save & New" resets form and maintains focus
- [x] Form submission calls `onsubmit` prop with correct data shape

---

## Phase 6: Keyboard Shortcuts & Final Polish - MOSTLY COMPLETED

### Deliverables

1. **Global Keyboard Shortcuts (when modal open)** - DONE
   - `⌘S` / `Ctrl+S`: Submit and close (prevent browser save dialog)
   - `⌘⇧S` / `Ctrl+Shift+S`: Submit and create another
   - `Escape`: Close modal (handled by shadcn Dialog)

2. **Cloze-Specific Shortcuts** - DONE
   - `Ctrl+1` through `Ctrl+9`: Insert cloze c1-c9 (wrap selection or insert template)
   - Only active when cloze textarea is focused (handled via `onkeydown` on textarea)

3. **Focus Management** - DONE
   - On modal open: focus first input of selected type (via `$effect`)
   - On type change: focus first input of new type (via `$effect`)
   - After "Save & New": focus first input (via `requestAnimationFrame`)
   - Focus trap: handled by shadcn Dialog

4. **ARIA Enhancements** - PARTIAL
   - [x] Modal: `role="dialog"`, `aria-modal="true"` (via shadcn Dialog)
   - [x] Type selector: `role="radiogroup"`, each option `role="radio"` with `aria-checked`
   - [ ] Form fields: `aria-describedby` for error messages, `aria-invalid` when invalid
   - [ ] Buttons: `aria-label` for icon-only buttons (mini-toolbar, cloze `+` button)

5. **Animation & Transitions** - DONE (via shadcn Dialog defaults)

6. **Edge Cases**
   - [x] Clicking backdrop closes modal
   - [x] Prevent body scroll when modal open (via shadcn Dialog)
   - [x] Handle rapid type switching gracefully
   - [ ] Prefill data populates correct fields and selects correct type (not implemented)

### Exit Criteria
- [x] All keyboard shortcuts work as specified
- [x] Focus management is correct in all scenarios
- [ ] Screen reader can navigate modal effectively (ARIA partially done)
- [x] Animations are smooth and non-distracting
- [ ] Prefill data works correctly (not implemented)
- [x] No console errors or warnings
- [ ] Modal works on mobile viewports (responsive) - needs testing

---

## Testing Checklist

### Functional Tests
- [ ] Create Basic fact with all fields populated
- [ ] Create Basic fact with only required fields
- [ ] Create Cloze fact with multiple cloze deletions
- [ ] Verify card count matches cloze deletion count
- [ ] "Save & New" creates fact and resets form
- [ ] Prefill data populates correctly
- [ ] Type switching preserves independent state
- [ ] All keyboard shortcuts function

### Accessibility Tests
- [ ] Navigate entire modal with keyboard only
- [ ] Screen reader announces modal, fields, and errors
- [ ] Focus is trapped and managed correctly
- [ ] Color contrast meets WCAG AA
- [ ] Touch targets are at least 44x44px on mobile

### Edge Cases
- [ ] Very long text in fields
- [ ] Rapid submission attempts
- [ ] Network error during submission
- [ ] Modal opened while another modal closing
- [ ] Browser back button while modal open

---

## Notes for Implementation

1. **Do not hardcode colors** - Use Tailwind design tokens and CSS custom properties from the existing theme.

2. **Follow existing patterns** - Look at other modal implementations in the codebase for consistent styling and behavior.

3. **TipTap integration is deferred** - mini-toolbar.svelte buttons should be visually complete but non-functional. Add TODO comments indicating where TipTap integration will occur.

4. **Image Occlusion is deferred** - Only implement the placeholder. Do not build any canvas, drawing, or region management functionality.

5. **Form actions vs. client submission** - This modal will likely use client-side submission via the proxy pattern mentioned in the tech stack. Confirm the preferred approach before Phase 5.

6. **State management** - Use Svelte 5 runes (`$state`, `$derived`) for all reactive state within components.

7. **Component composition** - Prefer small, focused components. Extract reusable pieces (icon-input.svelte, mini-toolbar.svelte) early.

8. **File naming convention** - Use lowercase with dashes for all component files (e.g., `create-fact-modal.svelte`, not `CreateFactModal.svelte`).

9. **Edit modal reusability** - Design editor components (basic-fact-editor.svelte, cloze-fact-editor.svelte) to accept initial values and work in both create and edit contexts. The edit modal implementation will follow this task.