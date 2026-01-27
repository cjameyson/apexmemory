## Prompt 1: API Types & Service Layer Foundation

**Goal:** Establish the type system and adapter pattern for notebooks.

### Deliverables

NOTE: frontend mockup used emoji and color, but these are not supported in the API.  If a user wants to add an emoji, they can simply include it as part of the name.  Color and other customizations will be added post MVP. `position` is also deferred til later and the default ordering will be alphabetical.  For now just use a 📘 icon for the emoji on the frontend to avoid changing all the components.

1. **Create `$lib/api/types.ts`**
   - `ApiNotebook` type (mirrors backend exactly: id, name, description, desired_retention, position, created_at, updated_at, archived_at)
   - `CreateNotebookRequest` type (name, description?)
   - `UpdateNotebookRequest` type (name?, description?, desired_retention?, position?)

2. **Update `$lib/types/notebook.ts`**
   - Keep existing UI fields (dueCount, totalCards, streak, retention)
   - Extend to include all `ApiNotebook` fields
   - Document which fields come from API vs computed/UI-only

3. **Create `$lib/services/notebooks.ts`**
   - `toNotebook(api: ApiNotebook): Notebook` adapter function
   - Merges mock stats (dueCount, totalCards, streak, retention) until stats API exists
   - Update `emoji` field to use 📘 icon on the frontend (will deprecate in a separate task)
   - localStorage helpers for UI preferences per notebook ID

### Validation Rules to Encode
- `name`: Required, 1-255 chars, trimmed
- `description`: Optional, max 10,000 chars
- `desired_retention`: Range (0.70-0.99), default 0.9
- User superforms and zod as needed

### Testing Checklist
- [ ] `toNotebook()` correctly transforms API response
- [ ] localStorage read/write works for UI preferences
- [ ] Types compile without errors

---

## Prompt 2: Layout Data Loading & Dropdown Navigation

**Goal:** Replace mock data with real API calls and simplify dropdown to navigation-only.


### Deliverables

1. **Update `(app)/+layout.server.ts`**
   - Fetch notebooks from API endpoint
   - Return notebooks in load function
   - Handle API errors gracefully (empty array fallback)

2. **Update `(app)/+layout.svelte`**
   - Receive notebooks from load function
   - Transform via `toNotebook()` adapter
   - Pass to dropdown component

3. **Update `$lib/components/navigation/notebooks-dropdown.svelte`**
   - **Remove:** Any edit/settings icons, inline actions
   - **Keep:** List of notebooks, due count badges, click-to-navigate
   - **Add:** Footer with "Create notebook" link (triggers event, modal comes later) and "Manage notebooks →" link to `/notebooks`
   - **Add:** Search filter when > 8 notebooks
   - **Add:** Empty state with friendly message and CTA

### Behavior Spec
- Lists active notebooks alphabetically
- Shows due count badge per notebook
- Click notebook → navigate to notebook dashboard
- No management actions in dropdown

### Testing Checklist
- [ ] Notebooks load from API on page load
- [ ] Dropdown displays notebooks alphabetically
- [ ] Due count badges show correctly
- [ ] Click navigates to correct notebook
- [ ] Empty state displays when no notebooks
- [ ] Search filter appears with > 8 notebooks
- [ ] API failure doesn't crash (graceful fallback)

---

## Prompt 3: Create Notebook Modal & Form Action

**Goal:** Enable quick notebook creation from anywhere in the app.

NOTE: Use superforms and zod as needed for consistent forms.

### Prerequisites
- Prompt 2 complete (dropdown has "Create notebook" trigger point)

### Deliverables

1. **Create `$lib/components/notebooks/create-notebook-modal.svelte`**
   - Modal overlay with form
   - Fields: Name (required), Description (optional)
   - Submit button with loading state
   - Cancel/close button
   - Inline validation errors

2. **Create or update form action for notebook creation**
   - POST to backend API
   - Handle validation errors (return to form)
   - Handle success (close modal, invalidate data)

3. **Wire up modal trigger**
   - "Create notebook" in dropdown footer opens modal
   - Modal can be opened from other places (event-based or store)

### UX Requirements
- Progressive enhancement (works without JS via form action)
- Loading state on submit button
- Inline error messages below fields
- Success: modal closes, notebook list refreshes
- Keyboard: Escape closes, Enter submits

### Not Included
- `desired_retention` field (use default 0.9)
- Any advanced settings

### Testing Checklist
- [ ] Modal opens from dropdown
- [ ] Empty name shows validation error
- [ ] Valid submission creates notebook
- [ ] Modal closes on success
- [ ] Notebook appears in dropdown after creation
- [ ] API error shows user-friendly message
- [ ] Escape key closes modal

---

## Prompt 4: Notebooks Management Page & Settings Page

**Goal:** Create dedicated pages for notebook management and configuration.

### Prerequisites
- Prompts 1-3 complete

### Deliverables

1. **Create `(app)/notebooks/+page.svelte` and `+page.server.ts`**
   - Header: "Notebooks" title + "Create notebook" button
   - Tabs: "Active" | "Archived"
   - Table/list view with columns: Name, Cards, Due, Retention, Created, Actions
   - Row click → navigate to notebook dashboard
   - Actions: Settings icon (→ settings page), Menu (archive)
   - Archived tab: same but with Unarchive action
   - Form actions: `archiveNotebook`, `unarchiveNotebook`

2. **Create `(app)/notebooks/[id]/settings/+page.svelte` and `+page.server.ts`**
   - Breadcrumb navigation
   - **General section:** Name input, Description textarea
   - **Spaced Repetition section:** Desired retention slider (0.70-0.99)
   - **Danger Zone section:** Archive button with confirmation dialog
   - Form action: `updateNotebook`
   - Archive redirects to `/notebooks`

3. **Update notebook dashboard component**
   - Add settings gear icon (⚙) in header
   - Click navigates to `/notebooks/[id]/settings`

### Form Actions Summary
| Route | Actions |
|-------|---------|
| `/notebooks` | `archiveNotebook`, `unarchiveNotebook` |
| `/notebooks/[id]/settings` | `updateNotebook`, `archiveNotebook` |

### Testing Checklist
- [ ] `/notebooks` lists all active notebooks
- [ ] Archived tab shows archived notebooks only
- [ ] Archive moves notebook to archived tab
- [ ] Unarchive moves notebook to active tab
- [ ] Settings page loads notebook data
- [ ] Name/description updates save correctly
- [ ] Retention slider constrains to 0.70-0.99
- [ ] Archive from settings redirects to /notebooks
- [ ] Dashboard has settings gear icon
- [ ] Gear icon navigates to settings page

---

## Summary

| Prompt | Focus | Key Files |
|--------|-------|-----------|
| 1 | Types & adapter pattern | `api/types.ts`, `types/notebook.ts`, `services/notebooks.ts` |
| 2 | API integration & dropdown | `+layout.server.ts`, `+layout.svelte`, `notebooks-dropdown.svelte` |
| 3 | Create modal & action | `create-notebook-modal.svelte`, form action |
| 4 | Management & settings pages | `/notebooks/*`, `/notebooks/[id]/settings/*` |

Each prompt is self-contained with clear inputs, outputs, and testing criteria. They should be executed in order since each builds on the previous.