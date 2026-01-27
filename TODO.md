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

### Backend
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

