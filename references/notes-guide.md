# Note and Card Processing

This document describes how notes and their associated cards are created, updated, and deleted in Apex Memory.

## Core Concepts

### Note to Card Relationship

A **note** is the user-editable content. A **card** is a reviewable unit with FSRS scheduling state.

| Note Type | Cards Generated | element_id Format |
|-----------|-----------------|-------------------|
| basic | 1 | `""` (empty string) |
| cloze | 1 per cloze deletion | `c1`, `c2`, ... `c999` |
| image_occlusion | 1 per mask region | `m_<nanoid>` |

### element_id Stability

The `element_id` column uniquely identifies which part of a note a card represents. These IDs are **stable**—deleting one element does not affect the IDs of others. This preserves review history for unchanged cards when notes are edited.

For cloze notes, element_ids are derived from the cloze markers in content (`{{c1::...}}`). For image occlusion, the client generates stable IDs (e.g., `m_abc123`) when creating mask regions.

## Note Creation

```
CreateNote(user_id, notebook_id, note_type, content):
    1. Validate content structure
    2. Extract element_ids from content
    3. Validate element count: 1 ≤ count ≤ 128
    
    4. BEGIN TRANSACTION
    5. INSERT note
    6. INSERT cards (one per element_id, all in 'new' state)
    7. COMMIT
    
    8. Return note with card count
```

### Element Extraction by Type

**Basic:**
```
element_ids = [""]
```

**Cloze:** Parse `{{c<N>::...}}` patterns from cloze_text fields.
```
Content: "The {{c1::mitochondria}} is the {{c2::powerhouse}}"
element_ids = ["c1", "c2"]
```

**Image Occlusion:** Read region IDs from content.
Read @./references/image-occlusion-note.md for full schema details.

```json
{
  "fields": [{
    "type": "image_occlusion",
    "regions": [
      {"id": "m_abc123", ...},
      {"id": "m_def456", ...}
    ]
  }]
}
```
```
element_ids = ["m_abc123", "m_def456"]
```

## Note Updates

Updates use a **diff algorithm** to minimize card churn and preserve review history.

```
UpdateNote(user_id, note_id, new_content):
    1. Fetch existing note
    2. Validate note_type unchanged (type changes not allowed)
    3. Validate new_content structure
    4. Extract expected element_ids from new_content
    5. Validate element count: 1 ≤ count ≤ 128
    
    6. BEGIN TRANSACTION
    
    7. Fetch existing element_ids:
       SELECT element_id FROM cards WHERE user_id = ? AND note_id = ?
    
    8. Compute diff:
       to_create = expected - existing
       to_delete = existing - expected
       unchanged = expected ∩ existing
    
    9. UPDATE note content
    
    10. DELETE cards WHERE element_id IN (to_delete)
        -- Reviews preserved with card_id = NULL
    
    11. INSERT new cards for to_create
        -- New cards start in 'new' state
    
    12. COMMIT
    
    13. Return { created, deleted, unchanged } counts
```

### Update Scenarios

| Scenario | Before | After | Action |
|----------|--------|-------|--------|
| Add cloze | `[c1]` | `[c1, c2]` | Create card for c2 |
| Remove cloze | `[c1, c2, c3]` | `[c1, c3]` | Delete card for c2 |
| Renumber cloze | `[c1, c2]` | `[c1, c4]` | Delete c2, create c4 |
| Content only | `[c1, c2]` | `[c1, c2]` | Update note, no card changes |
| Add mask region | `[m_a, m_b]` | `[m_a, m_b, m_c]` | Create card for m_c |
| Remove mask region | `[m_a, m_b, m_c]` | `[m_a, m_c]` | Delete card for m_b |

## Note Deletion

When a note is deleted:
1. All associated cards are deleted (FK cascade)
2. Reviews are preserved with `card_id = NULL` and `note_id = NULL`
3. Review data remains available for optimizer training and historical stats

## Validation Rules

### Content Structure
```json
{
  "version": 1,
  "fields": [...]
}
```
- `version` must be a number
- `fields` must be an array

### element_id Format
| Type | Pattern | Examples |
|------|---------|----------|
| basic | `^$` | `""` |
| cloze | `^c[1-9][0-9]{0,2}$` | `c1`, `c42`, `c999` |
| image_occlusion | `^m_[a-zA-Z0-9_-]{6,24}$` | `m_abc123`, `m_xK9-mQ_2` |

### Constraints
- **Max cards per note:** 128
- **Note type immutable:** Cannot change after creation
- **At least one element:** Notes must generate ≥1 card

### Cloze-Specific Rules
- Gaps allowed: `[c1, c3]` is valid (c2 skipped)
- Duplicates collapsed: `{{c1::a}} {{c1::b}}` → single card showing both blanks
- Zero not allowed: `c0` is invalid (1-indexed)

## Review History Preservation

When cards are deleted (directly or via note update), their reviews are **not deleted**. Instead:

| Column | Behavior |
|--------|----------|
| `card_id` | Set to NULL |
| `note_id` | Set to NULL (if note deleted) |
| `notebook_id` | Preserved (was already denormalized) |
| All review data | Preserved |

This allows:
- Accurate lifetime statistics
- FSRS optimizer training on full history
- Analytics even after content cleanup

## Error Handling

| Error | Condition |
|-------|-----------|
| `ValidationError` | Invalid content structure |
| `ValidationError` | Note type change attempted |
| `ValidationError` | Zero elements extracted |
| `ValidationError` | More than 128 elements |
| `ValidationError` | Invalid element_id format |
| `NotFoundError` | Note doesn't exist or wrong user |