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

## Phase 1: Modal Shell & Type Selector

### Deliverables

1. **create-fact-modal.svelte**
   - Modal overlay with backdrop blur
   - Fixed positioning: 40px from viewport top, centered horizontally
   - Max height: `calc(100vh - 80px)` with internal scroll
   - Subtle border radius (use design token for `rounded-lg` equivalent)
   - Three sections: Header, Content (scrollable), Footer
   
2. **Modal Header**
   - Compact height (~48px)
   - Title: "Create Fact"
   - Close button (X icon) on right
   
3. **fact-type-selector.svelte**
   - Three options displayed horizontally as radio-button style cards
   - Each option shows: icon, label, short description
   - Selected state: accent border, tinted background, checkmark
   - Options:
     - Basic: "Front & back flashcard"
     - Cloze: "Fill in the blanks"  
     - Image Occlusion: "Hide parts of an image"
   
4. **Modal Footer**
   - Left: Keyboard shortcut hints (`⌘S` save, `⌘⇧S` save & new, `Esc` cancel)
   - Right: Card count indicator, Cancel button, "Save & New" button, "Create Fact" primary button

5. **Props Interface**
   ```typescript
   interface CreateFactModalProps {
     open: boolean;
     notebookId: string;
     prefillData?: {
       type?: 'basic' | 'cloze' | 'image_occlusion';
       front?: string;
       back?: string;
       clozeText?: string;
     };
     onClose: () => void;
     onSubmit: (data: FactFormData) => Promise<void>;
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

### Human Review Checkpoint
Verify modal positioning, type selector interaction, and basic accessibility before proceeding.

---

## Phase 2: Basic Fact Editor - COMPLETED

### Deliverables

1. **basic-fact-editor.svelte**
   - Four fields in order: Front, Back, Back Extra, Hint
   - Props: `initialData` (optional, for edit mode), `onchange` callback
   - Design for reuse in edit modal (will receive existing fact data)
   
2. **Front Field**
   - Label: "Front"
   - mini-toolbar aligned right of label (placeholder buttons for now)
   - Textarea with border, 100px min-height
   - Placeholder: "Enter the question or prompt..."
   - Auto-focus on modal open when Basic type selected
   
3. **Back Field**
   - Label: "Back"  
   - mini-toolbar aligned right of label
   - Textarea with border, 100px min-height
   - Placeholder: "Enter the answer..."
   
4. **Back Extra Field** (using icon-input)
   - Embedded Info icon on left inside input
   - Multiline textarea, ~80px min-height
   - Placeholder: "Extra information shown with the answer (optional)"
   - No separate label (icon serves as visual indicator)
   
5. **Hint Field** (using icon-input)
   - Embedded Lightbulb icon on left inside input
   - Single-line input
   - Placeholder: "Hint to help recall the answer (optional)"
   - No separate label

6. **icon-input.svelte** (reusable)
   - Props: `icon`, `placeholder`, `value`, `multiline`, `onChange`
   - Icon positioned absolutely inside input padding

7. **mini-toolbar.svelte** (placeholder)
   - Row of icon buttons: Bold, Italic, Underline | Code, Link | Bullet List, Numbered List
   - Buttons are non-functional placeholders (will integrate TipTap later)
   - Appropriate spacing and dividers between groups

8. **State Management**
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

## Phase 3: Cloze Fact Editor

### Deliverables

1. **cloze-fact-editor.svelte**
   - Two fields: Cloze Text, Back Extra
   - Props: `initialData` (optional, for edit mode), `onchange` callback
   - Design for reuse in edit modal (will receive existing fact data)
   
2. **Cloze Text Field**
   - Label: "Cloze Text"
   - Toolbar row containing:
     - mini-toolbar (left)
     - Divider
     - "Insert:" label
     - Cloze buttons: `c1`, `c2`, `c3` + `+` button for next available
   - Monospace font textarea, 160px min-height
   - Placeholder text explaining syntax with example
   - Helper text below: keyboard shortcut hint for Ctrl+1-9
   
3. **Cloze Insertion Logic**
   - Track which cloze numbers are used in current text
   - Used numbers: show button in muted/secondary style
   - Unused numbers: show button in accent/primary style
   - `+` button: inserts next unused cloze number
   - Insertion behavior:
     - If text selected: wrap selection in `{{cN::selection}}`
     - If no selection: insert `{{cN::}}` and place cursor between `::`  and `}}`
   - After insertion, return focus to textarea at appropriate cursor position
   
4. **Back Extra Field**
   - Same icon-input as Basic editor
   - Placeholder: "Extra information shown with each answer (optional)"

5. **Card Count Calculation**
   - Parse cloze text for unique cloze numbers
   - Display count: "{N} cards" in footer
   - Regex: `/\{\{c(\d+)::/g` to find cloze markers

### Exit Criteria
- [ ] Cloze text field has monospace font
- [ ] Cloze buttons show correct used/unused states
- [ ] Clicking cloze button wraps selected text correctly
- [ ] Clicking cloze button with no selection inserts template
- [ ] Cursor position is correct after insertion
- [ ] `+` button inserts next available number
- [ ] Card count updates dynamically as cloze markers are added/removed
- [ ] Back Extra field works identically to Basic editor
- [ ] Tab order is logical

### Human Review Checkpoint
Verify cloze insertion behavior and card count calculation before proceeding.

---

## Phase 4: Image Occlusion Placeholder

### Deliverables

1. **image-occlusion-placeholder.svelte**
   - Centered content area with:
     - Large icon (Image or similar)
     - Heading: "Image Occlusion"
     - Subheading: "Coming Soon"
     - Brief description: "Create flashcards by hiding parts of diagrams, charts, and images."
   - Muted/disabled visual treatment
   - Fills the content area of the modal

2. **Disabled State**
   - When Image Occlusion is selected, footer buttons remain but "Create Fact" is disabled
   - Card count shows "0 cards"
   - Tooltip or visual indicator that this feature is not yet available

### Exit Criteria
- [ ] Placeholder displays when Image Occlusion type selected
- [ ] Create Fact button is disabled
- [ ] Visual design is consistent with overall modal
- [ ] User understands this is a future feature

### Human Review Checkpoint
Verify placeholder communicates "coming soon" clearly before proceeding.

---

## Phase 5: Form Validation & Submission

### Deliverables

1. **Zod Schemas**
   ```typescript
   const basicFactSchema = z.object({
     factType: z.literal('basic'),
     front: z.string().min(1, 'Front is required'),
     back: z.string().min(1, 'Back is required'),
     backExtra: z.string().optional(),
     hint: z.string().optional(),
   });
   
   const clozeFactSchema = z.object({
     factType: z.literal('cloze'),
     clozeText: z.string()
       .min(1, 'Cloze text is required')
       .refine(
         (val) => /\{\{c\d+::.+?\}\}/.test(val),
         'At least one cloze deletion is required'
       ),
     backExtra: z.string().optional(),
   });
   
   const factFormSchema = z.discriminatedUnion('factType', [
     basicFactSchema,
     clozeFactSchema,
   ]);
   ```

2. **Validation Behavior**
   - Validate on submit attempt
   - Show inline error messages below invalid fields
   - Focus first invalid field on validation failure
   - Clear field error when user starts typing

3. **Submit Flow**
   - "Create Fact" button: validate → submit → close modal
   - "Save & New" button: validate → submit → reset form (stay on same type)
   - Show loading state on buttons during submission
   - Handle submission errors gracefully (toast or inline message)

4. **Form Reset**
   - After "Save & New", clear only the active type's fields
   - Preserve the selected fact type
   - Return focus to first field

### Exit Criteria
- [ ] Empty required fields show validation errors
- [ ] Cloze text without any cloze markers shows error
- [ ] Errors clear when user edits field
- [ ] Submit shows loading state
- [ ] "Save & New" resets form and maintains focus
- [ ] Form submission calls `onSubmit` prop with correct data shape

### Human Review Checkpoint
Verify validation messages and submission flow before proceeding.

---

## Phase 6: Keyboard Shortcuts & Final Polish

### Deliverables

1. **Global Keyboard Shortcuts (when modal open)**
   - `⌘S` / `Ctrl+S`: Submit and close (prevent browser save dialog)
   - `⌘⇧S` / `Ctrl+Shift+S`: Submit and create another
   - `Escape`: Close modal (if no unsaved changes, or confirm)
   
2. **Cloze-Specific Shortcuts**
   - `Ctrl+1` through `Ctrl+9`: Insert cloze c1-c9 (wrap selection or insert template)
   - Only active when cloze textarea is focused

3. **Focus Management**
   - On modal open: focus first input of selected type
   - On type change: focus first input of new type
   - After "Save & New": focus first input
   - Focus trap: Tab/Shift+Tab cycle within modal

4. **ARIA Enhancements**
   - Modal: `role="dialog"`, `aria-modal="true"`, `aria-labelledby` pointing to title
   - Type selector: `role="radiogroup"`, each option `role="radio"` with `aria-checked`
   - Form fields: `aria-describedby` for error messages, `aria-invalid` when invalid
   - Buttons: clear `aria-label` where icon-only

5. **Animation & Transitions**
   - Modal entrance: fade in backdrop, scale up modal slightly
   - Modal exit: reverse of entrance
   - Type selector: smooth background/border transitions
   - Keep animations subtle and fast (150-200ms)

6. **Edge Cases**
   - Clicking backdrop closes modal
   - Prevent body scroll when modal open
   - Handle rapid type switching gracefully
   - Prefill data populates correct fields and selects correct type

### Exit Criteria
- [ ] All keyboard shortcuts work as specified
- [ ] Focus management is correct in all scenarios
- [ ] Screen reader can navigate modal effectively
- [ ] Animations are smooth and non-distracting
- [ ] Prefill data works correctly
- [ ] No console errors or warnings
- [ ] Modal works on mobile viewports (responsive)

### Human Review Checkpoint
Final review of complete component before integration.

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