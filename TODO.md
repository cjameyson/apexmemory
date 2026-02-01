# TODO

# TODO

## Ideas

### Source Decomposition
Large PDFs (up to 1GB, 500 pages) should be broken into logical chunks (chapters, sections) for easier processing and chat navigation. Detection via PDF inspection, OCR, and LLM. Requires a sophisticated source processing pipeline.

### Deduplication
Sources are commonly shared across users at a school (e.g., nursing students with same lecture PDF/video). Need strategy to avoid duplicate processing of large files, OCR, and embedding.

### Source Viewers / Content Components
Automatic card recommendation based on current card set and source content. Find potential gaps in knowledge and suggest cards to fill those gaps.

**PDF and Document Based Sources**
- PDF displayed in left pane, cards/summary/chat in right pane
- Select text -> take action; select image -> take action
- Cards tab updates as user scrolls to show linked cards (advanced)

**Audio**
- Audio player with optional transcript
- Select text in transcript -> take action (capture timestamp)
- Cards tab updates as user scrolls (advanced)

**Video**
- Video player with optional transcript
- Select text in transcript -> take action (capture timestamp)
- Cards tab updates as user watches to show linked cards (advanced)

### Practice / Past Tests
Add a special source type for practice tests or previous exams. This allows the AI system to process them specifically to:
- Extract key insights on test problems and patterns.
- Compare test content to existing flashcards to identify gaps.

---

## Active

### Bugs

### Frontend
- [ ] Generic 'select text' -> take action capability (create cards, add to notes, summarize, explain)
- [ ] Screen grab capability for image occlusion cards (or support paste from native screenshot tool)
- [ ] **image_occlusion region counting may be incorrect** - `fact-display.ts` line 61 filters fields by `f.type === 'image_occlusion'` to count regions. Verify against actual image_occlusion content shape from backend (the regions are nested inside the field, not separate fields). The current logic likely returns 1 (the image field itself) rather than the actual region count.
- [ ] **`stripHtml` is naive** - `fact-display.ts` uses `/<[^>]*>/g` which breaks on edge cases like `<img alt="a>b">` or malformed HTML. Acceptable for display-only truncation but consider a proper sanitizer if user-authored HTML is allowed in fields.
- [ ] **Expanded row `<tr>` has no keyed association** - In `FactTableRow.svelte`, the `{#if expanded}` block renders a second `<tr>` outside the `{#each}` key scope. If Svelte reorders rows during keyed updates, the expanded row could theoretically detach. Not a practical issue with current structure but worth noting.


### Backend
- [ ] Slugify notebook and sources on creation/edit - append nanoid to end if collision.  Use slugs for URLs
- [ ] Determine best way to automatically apply migrations to the test database during unit tests
- [ ] Add pagination to ListNotebooks endpoint (deferred: max notebooks per user is sub-100)
- [ ] Consider storing desired_retention as separate column for performance optimization on high-volume lists
- [ ] Add PUT /v1/notebooks/{id}/unarchive endpoint (query exists but not exposed)

---

## Done
- [x] Use middleware chaining library like 'alice' for cleaner middleware setup and reduce auth check boilerplate in handlers
- [x] Source detail refactor: two-pane layout with collapsible right sidebar
- [x] Icon-only buttons always have tooltips
- [x] Use lucide-maximize2/minimize2 for expand/collapse buttons
- [x] Update source tabs to pill-based styling
- [x] Improve source tab styling
- [x] Main pane expand remembers side panel state
- [x] When clicking the Review dropdown in the topnav, we enter focus mode, but if I press space to start a review, the dropdown is shown again ontop of the review interface.

