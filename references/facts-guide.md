# Fact and Card Processing

This document describes how facts and their associated cards are created, updated, and deleted in Apex Memory.

## Core Concepts

### Fact to Card Relationship

A **fact** is the user-editable content. A **card** is a reviewable unit with FSRS scheduling state.

| Fact Type | Cards Generated | element_id Format |
|-----------|-----------------|-------------------|
| basic | 1 | `""` (empty string) |
| cloze | 1 per cloze deletion | `c1`, `c2`, ... `c999` |
| image_occlusion | 1 per mask region | `m_<nanoid>` |

### element_id Stability

The `element_id` column uniquely identifies which part of a fact a card represents. These IDs are **stable**—deleting one element does not affect the IDs of others. This preserves review history for unchanged cards when facts are edited.

For cloze facts, element_ids are derived from the cloze markers in content (`{{c1::...}}`). For image occlusion, the client generates stable IDs (e.g., `m_abc123`) when creating mask regions.

## Fact Creation

```
CreateFact(user_id, notebook_id, fact_type, content):
    1. Validate content structure
    2. Extract element_ids from content
    3. Validate element count: 1 ≤ count ≤ 128
    
    4. BEGIN TRANSACTION
    5. INSERT fact
    6. INSERT cards (one per element_id, all in 'new' state)
    7. COMMIT
    
    8. Return fact with card count
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
Read @./references/image-occlusion-fact-type.md for full schema details.

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

## Fact Updates

Updates use a **diff algorithm** to minimize card churn and preserve review history.

```
UpdateFact(user_id, fact_id, new_content):
    1. Fetch existing fact
    2. Validate fact_type unchanged (type changes not allowed)
    3. Validate new_content structure
    4. Extract expected element_ids from new_content
    5. Validate element count: 1 ≤ count ≤ 128
    
    6. BEGIN TRANSACTION
    
    7. Fetch existing element_ids:
       SELECT element_id FROM cards WHERE user_id = ? AND fact_id = ?
    
    8. Compute diff:
       to_create = expected - existing
       to_delete = existing - expected
       unchanged = expected ∩ existing
    
    9. UPDATE fact content
    
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
| Content only | `[c1, c2]` | `[c1, c2]` | Update fact, no card changes |
| Add mask region | `[m_a, m_b]` | `[m_a, m_b, m_c]` | Create card for m_c |
| Remove mask region | `[m_a, m_b, m_c]` | `[m_a, m_c]` | Delete card for m_b |

## Fact Deletion

When a fact is deleted:
1. All associated cards are deleted (FK cascade)
2. Reviews are preserved with `card_id = NULL` and `fact_id = NULL`
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
- **Max cards per fact:** 128
- **Fact type immutable:** Cannot change after creation
- **At least one element:** Facts must generate ≥1 card

### Cloze-Specific Rules
- Gaps allowed: `[c1, c3]` is valid (c2 skipped)
- Duplicates collapsed: `{{c1::a}} {{c1::b}}` → single card showing both blanks
- Zero not allowed: `c0` is invalid (1-indexed)

## Review History Preservation

When cards are deleted (directly or via fact update), their reviews are **not deleted**. Instead:

| Column | Behavior |
|--------|----------|
| `card_id` | Set to NULL |
| `fact_id` | Set to NULL (if fact deleted) |
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
| `ValidationError` | Fact type change attempted |
| `ValidationError` | Zero elements extracted |
| `ValidationError` | More than 128 elements |
| `ValidationError` | Invalid element_id format |
| `NotFoundError` | Fact doesn't exist or wrong user |